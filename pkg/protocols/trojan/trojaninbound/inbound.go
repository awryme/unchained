package trojaninbound

import (
	"context"
	"net"
	"os"

	"github.com/awryme/unchained/pkg/protocols/trojan/trojanproto"
	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/adapter/inbound"
	"github.com/sagernet/sing-box/common/listener"
	"github.com/sagernet/sing-box/common/mux"
	"github.com/sagernet/sing-box/common/tls"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/transport/v2ray"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/auth"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

const TypeTrojanV2 = "trojan_v2"

func Register(registry *inbound.Registry) {
	inbound.Register(registry, TypeTrojanV2, NewInbound)
}

var _ adapter.TCPInjectableInbound = (*Inbound)(nil)

type Inbound struct {
	inbound.Adapter
	router                   adapter.ConnectionRouterEx
	logger                   log.ContextLogger
	listener                 *listener.Listener
	service                  *trojanproto.Service
	tlsConfig                tls.ServerConfig
	fallbackAddr             M.Socksaddr
	fallbackAddrTLSNextProto map[string]M.Socksaddr
	transport                adapter.V2RayServerTransport
}

func NewInbound(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options TrojanInboundOptions) (adapter.Inbound, error) {
	inbound := &Inbound{
		Adapter: inbound.NewAdapter(C.TypeTrojan, tag),
		router:  router,
		logger:  logger,
	}
	if options.TLS != nil {
		tlsConfig, err := tls.NewServer(ctx, logger, common.PtrValueOrDefault(options.TLS))
		if err != nil {
			return nil, err
		}
		inbound.tlsConfig = tlsConfig
	}
	var fallbackHandler N.TCPConnectionHandlerEx
	if options.Fallback != nil && options.Fallback.Server != "" || len(options.FallbackForALPN) > 0 {
		if options.Fallback != nil && options.Fallback.Server != "" {
			inbound.fallbackAddr = options.Fallback.Build()
			if !inbound.fallbackAddr.IsValid() {
				return nil, E.New("invalid fallback address: ", inbound.fallbackAddr)
			}
		}
		if len(options.FallbackForALPN) > 0 {
			if inbound.tlsConfig == nil {
				return nil, E.New("fallback for ALPN is not supported without TLS")
			}
			fallbackAddrNextProto := make(map[string]M.Socksaddr)
			for nextProto, destination := range options.FallbackForALPN {
				fallbackAddr := destination.Build()
				if !fallbackAddr.IsValid() {
					return nil, E.New("invalid fallback address for ALPN ", nextProto, ": ", fallbackAddr)
				}
				fallbackAddrNextProto[nextProto] = fallbackAddr
			}
			inbound.fallbackAddrTLSNextProto = fallbackAddrNextProto
		}
		fallbackHandler = adapter.NewUpstreamContextHandlerEx(inbound.fallbackConnection, nil)
	}
	service := trojanproto.NewService(
		adapter.NewUpstreamContextHandlerEx(inbound.newConnection, inbound.newPacketConnection),
		fallbackHandler,
		logger,
		options.UserStore,
	)
	var err error
	if options.Transport != nil {
		inbound.transport, err = v2ray.NewServerTransport(ctx, logger, common.PtrValueOrDefault(options.Transport), inbound.tlsConfig, (*inboundTransportHandler)(inbound))
		if err != nil {
			return nil, E.Cause(err, "create server transport: ", options.Transport.Type)
		}
	}
	inbound.router, err = mux.NewRouterWithOptions(inbound.router, logger, common.PtrValueOrDefault(options.Multiplex))
	if err != nil {
		return nil, err
	}
	inbound.service = service
	inbound.listener = listener.New(listener.Options{
		Context:           ctx,
		Logger:            logger,
		Network:           []string{N.NetworkTCP},
		Listen:            options.ListenOptions,
		ConnectionHandler: inbound,
	})
	return inbound, nil
}

func (h *Inbound) Start(stage adapter.StartStage) error {
	if stage != adapter.StartStateStart {
		return nil
	}
	if h.tlsConfig != nil {
		err := h.tlsConfig.Start()
		if err != nil {
			return E.Cause(err, "create TLS config")
		}
	}
	if h.transport == nil {
		return h.listener.Start()
	}
	if common.Contains(h.transport.Network(), N.NetworkTCP) {
		tcpListener, err := h.listener.ListenTCP()
		if err != nil {
			return err
		}
		go func() {
			sErr := h.transport.Serve(tcpListener)
			if sErr != nil && !E.IsClosed(sErr) {
				h.logger.Error("transport serve error: ", sErr)
			}
		}()
	}
	if common.Contains(h.transport.Network(), N.NetworkUDP) {
		udpConn, err := h.listener.ListenUDP()
		if err != nil {
			return err
		}
		go func() {
			sErr := h.transport.ServePacket(udpConn)
			if sErr != nil && !E.IsClosed(sErr) {
				h.logger.Error("transport serve error: ", sErr)
			}
		}()
	}
	return nil
}

func (h *Inbound) Close() error {
	return common.Close(
		h.listener,
		h.tlsConfig,
		h.transport,
	)
}

func (h *Inbound) NewConnectionEx(ctx context.Context, conn net.Conn, metadata adapter.InboundContext, onClose N.CloseHandlerFunc) {
	if h.tlsConfig != nil && h.transport == nil {
		tlsConn, err := tls.ServerHandshake(ctx, conn, h.tlsConfig)
		if err != nil {
			N.CloseOnHandshakeFailure(conn, onClose, err)
			h.logger.ErrorContext(ctx, E.Cause(err, "process connection from ", metadata.Source, ": TLS handshake"))
			return
		}
		conn = tlsConn
	}
	err := h.service.NewConnection(adapter.WithContext(ctx, &metadata), conn, metadata.Source, onClose)
	if err != nil {
		N.CloseOnHandshakeFailure(conn, onClose, err)
		h.logger.ErrorContext(ctx, E.Cause(err, "process connection from ", metadata.Source))
	}
}

func (h *Inbound) newConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext, onClose N.CloseHandlerFunc) {
	metadata.Inbound = h.Tag()
	metadata.InboundType = h.Type()
	user, loaded := auth.UserFromContext[trojanproto.User](ctx)
	if !loaded {
		N.CloseOnHandshakeFailure(conn, onClose, os.ErrInvalid)
		return
	}
	metadata.User = user.Name
	h.logger.InfoContext(ctx, "[", user.Name, "] inbound connection to ", metadata.Destination)
	h.router.RouteConnectionEx(ctx, conn, metadata, onClose)
}

func (h *Inbound) newPacketConnection(ctx context.Context, conn N.PacketConn, metadata adapter.InboundContext, onClose N.CloseHandlerFunc) {
	metadata.Inbound = h.Tag()
	metadata.InboundType = h.Type()
	user, loaded := auth.UserFromContext[trojanproto.User](ctx)
	if !loaded {
		N.CloseOnHandshakeFailure(conn, onClose, os.ErrInvalid)
		return
	}
	metadata.User = user.Name
	h.logger.InfoContext(ctx, "[", user.Name, "] inbound packet connection to ", metadata.Destination)
	h.router.RoutePacketConnectionEx(ctx, conn, metadata, onClose)
}

func (h *Inbound) fallbackConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext, onClose N.CloseHandlerFunc) {
	var fallbackAddr M.Socksaddr
	if len(h.fallbackAddrTLSNextProto) > 0 {
		if tlsConn, loaded := common.Cast[tls.Conn](conn); loaded {
			connectionState := tlsConn.ConnectionState()
			if connectionState.NegotiatedProtocol != "" {
				if fallbackAddr, loaded = h.fallbackAddrTLSNextProto[connectionState.NegotiatedProtocol]; !loaded {
					h.logger.DebugContext(ctx, "process connection from ", metadata.Source, ": fallback disabled for ALPN: ", connectionState.NegotiatedProtocol)
					N.CloseOnHandshakeFailure(conn, onClose, os.ErrInvalid)
					return
				}
			}
		}
	}
	if !fallbackAddr.IsValid() {
		if !h.fallbackAddr.IsValid() {
			h.logger.DebugContext(ctx, "process connection from ", metadata.Source, ": fallback disabled by default")
			N.CloseOnHandshakeFailure(conn, onClose, os.ErrInvalid)
			return
		}
		fallbackAddr = h.fallbackAddr
	}
	metadata.Inbound = h.Tag()
	metadata.InboundType = h.Type()
	metadata.Destination = fallbackAddr
	h.logger.InfoContext(ctx, "fallback connection to ", fallbackAddr)
	h.router.RouteConnectionEx(ctx, conn, metadata, onClose)
}

var _ adapter.V2RayServerTransportHandler = (*inboundTransportHandler)(nil)

type inboundTransportHandler Inbound

func (h *inboundTransportHandler) NewConnectionEx(ctx context.Context, conn net.Conn, source M.Socksaddr, destination M.Socksaddr, onClose N.CloseHandlerFunc) {
	var metadata adapter.InboundContext
	metadata.Source = source
	metadata.Destination = destination
	//nolint:staticcheck
	metadata.InboundDetour = h.listener.ListenOptions().Detour
	//nolint:staticcheck
	metadata.InboundOptions = h.listener.ListenOptions().InboundOptions
	h.logger.InfoContext(ctx, "inbound connection from ", metadata.Source)
	(*Inbound)(h).NewConnectionEx(ctx, conn, metadata, onClose)
}
