package singboxserver

import (
	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/constants"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

func makeTlsOptions(cfg appconfig.Reality) option.InboundTLSOptionsContainer {
	return option.InboundTLSOptionsContainer{
		TLS: &option.InboundTLSOptions{
			Enabled:    true,
			ServerName: cfg.Server,
			Reality: &option.InboundRealityOptions{
				Enabled: true,
				Handshake: option.InboundRealityHandshakeOptions{
					ServerOptions: option.ServerOptions{
						Server:     cfg.Server,
						ServerPort: constants.DefaultRealityServerPort,
					},
				},
				PrivateKey: cfg.PrivateKey,
				ShortID: badoption.Listable[string]{
					cfg.ShortId,
				},
				MaxTimeDifference: badoption.Duration(cfg.TimeDiff),
			},
		},
	}
}
