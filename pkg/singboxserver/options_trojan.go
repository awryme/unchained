package singboxserver

import (
	"net/netip"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/protocols"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

func makeTrojanInboundOptions(cfg appconfig.Config, addr netip.Addr) option.Inbound {
	return option.Inbound{
		Type: protocols.Trojan,
		Tag:  "inbound-" + protocols.Trojan,
		Options: &option.TrojanInboundOptions{
			ListenOptions: option.ListenOptions{
				Listen:     (*badoption.Addr)(&addr),
				ListenPort: uint16(cfg.Listen.Port),
			},
			Users: []option.TrojanUser{
				{
					Name:     cfg.Name(),
					Password: cfg.TrojanPassword,
				},
			},
			InboundTLSOptionsContainer: makeTlsOptions(cfg.Reality),
		},
	}
}
