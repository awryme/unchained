package singboxserver

import (
	"net/netip"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/protocols"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

const vlessFlowVision = "xtls-rprx-vision"

func makeVlessInboundOptions(cfg appconfig.Config, addr netip.Addr) option.Inbound {
	return option.Inbound{
		Type: protocols.Vless,
		Tag:  "inbound-" + protocols.Vless,
		Options: &option.VLESSInboundOptions{
			ListenOptions: option.ListenOptions{
				Listen:     (*badoption.Addr)(&addr),
				ListenPort: uint16(cfg.Listen.Port),
			},
			Users: []option.VLESSUser{
				{
					Name: cfg.Name(),
					UUID: cfg.VlessUUID,
					Flow: vlessFlowVision,
				},
			},
			InboundTLSOptionsContainer: makeTlsOptions(cfg.Reality),
		},
	}
}
