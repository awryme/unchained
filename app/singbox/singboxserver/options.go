package singboxserver

import (
	"net/netip"

	"github.com/awryme/unchained/app/appconfig"
	"github.com/awryme/unchained/pkg/reality"
	"github.com/sagernet/sing-box/option"
	dns "github.com/sagernet/sing-dns"
	"github.com/sagernet/sing/common/json/badoption"
)

func makeLogOptions(cfg appconfig.Singbox) *option.LogOptions {
	return &option.LogOptions{
		Disabled:     false,
		Timestamp:    true,
		DisableColor: true,

		Level: cfg.LogLevel,
	}
}

func makeDnsOptions(cfg appconfig.Singbox) *option.DNSOptions {
	server := option.DNSServerOptions{
		Tag:     "dns-out",
		Address: cfg.DNS,
	}
	if cfg.DNSIPv4Only {
		server.Strategy = option.DomainStrategy(dns.DomainStrategyUseIPv4)
	}
	return &option.DNSOptions{
		Servers: []option.DNSServerOptions{server},
	}
}

func makeListenOptions(listen netip.AddrPort) option.ListenOptions {
	addr := listen.Addr()
	return option.ListenOptions{
		Listen:     (*badoption.Addr)(&addr),
		ListenPort: uint16(listen.Port()),
	}
}

func makeTlsOptions(cfg appconfig.Reality) option.InboundTLSOptions {
	return option.InboundTLSOptions{
		Enabled:    true,
		ServerName: cfg.Server,
		Reality: &option.InboundRealityOptions{
			Enabled: true,
			Handshake: option.InboundRealityHandshakeOptions{
				ServerOptions: option.ServerOptions{
					Server:     cfg.Server,
					ServerPort: reality.DefaultServerPort,
				},
			},
			PrivateKey: cfg.PrivateKey,
			ShortID: badoption.Listable[string]{
				cfg.ShortId,
			},
			MaxTimeDifference: badoption.Duration(cfg.TimeDiff),
		},
	}
}
