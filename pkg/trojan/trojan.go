package trojan

import (
	"context"
	"fmt"
	"io"
	"net/netip"

	"github.com/awryme/unchained/appconfig"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

func Run(ctx context.Context, cfg appconfig.Config) (io.Closer, error) {
	addr, err := netip.ParseAddr(cfg.ListenAddr)
	if err != nil {
		return nil, fmt.Errorf("parse ip addr to listen on: %w", err)
	}
	ctx = box.Context(ctx,
		include.InboundRegistry(),
		include.OutboundRegistry(),
		include.EndpointRegistry(),
	)

	instance, err := box.New(box.Options{
		Context: ctx,
		Options: option.Options{
			Log: makeLogOptions(cfg),
			DNS: makeDnsOptions(cfg),
			Inbounds: []option.Inbound{
				makeTrojanInboundOptions(cfg, addr),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("make sing-box instance: %w", err)
	}
	err = instance.Start()
	if err != nil {
		return nil, fmt.Errorf("start sing-box instance: %w", err)
	}
	return instance, nil
}

func makeLogOptions(cfg appconfig.Config) *option.LogOptions {
	return &option.LogOptions{
		Disabled:     false,
		Timestamp:    true,
		DisableColor: true,

		Level: cfg.LogLevel,
	}
}

func makeDnsOptions(cfg appconfig.Config) *option.DNSOptions {
	return &option.DNSOptions{
		Servers: []option.DNSServerOptions{
			{
				Tag:     "dns-out",
				Address: cfg.DNS,
			},
		},
	}
}

func makeTrojanInboundOptions(cfg appconfig.Config, addr netip.Addr) option.Inbound {
	return option.Inbound{
		Type: "trojan",
		Tag:  "inbound-trojan",
		Options: &option.TrojanInboundOptions{
			ListenOptions: option.ListenOptions{
				Listen:     (*badoption.Addr)(&addr),
				ListenPort: uint16(cfg.Port),
			},
			Users: []option.TrojanUser{
				{
					Name:     "app",
					Password: cfg.Password,
				},
			},
			InboundTLSOptionsContainer: option.InboundTLSOptionsContainer{
				TLS: &option.InboundTLSOptions{
					Enabled: true,
					Reality: &option.InboundRealityOptions{
						Enabled: true,
						Handshake: option.InboundRealityHandshakeOptions{
							ServerOptions: option.ServerOptions{
								Server:     cfg.RealityServer,
								ServerPort: 443,
							},
						},
						PrivateKey: cfg.RealityPrivateKey,
						ShortID: badoption.Listable[string]{
							cfg.RealityShortId,
						},
						MaxTimeDifference: badoption.Duration(cfg.RealityTimeDiff),
					},
				},
			},
		},
	}
}
