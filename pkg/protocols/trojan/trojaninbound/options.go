package trojaninbound

import (
	"github.com/awryme/unchained/pkg/protocols/trojan/trojanproto"
	"github.com/sagernet/sing-box/option"
)

type TrojanInboundOptions struct {
	option.ListenOptions
	UserStore trojanproto.UserStore
	option.InboundTLSOptionsContainer
	Fallback        *option.ServerOptions            `json:"fallback,omitempty"`
	FallbackForALPN map[string]*option.ServerOptions `json:"fallback_for_alpn,omitempty"`
	Multiplex       *option.InboundMultiplexOptions  `json:"multiplex,omitempty"`
	Transport       *option.V2RayTransportOptions    `json:"transport,omitempty"`
}

func MakeInbound(userStore trojanproto.UserStore, listen option.ListenOptions, tls option.InboundTLSOptions) option.Inbound {
	return option.Inbound{
		Type: TypeTrojanV2,
		Tag:  "inbound-" + TypeTrojanV2,
		Options: &TrojanInboundOptions{
			ListenOptions: listen,
			UserStore:     userStore,
			InboundTLSOptionsContainer: option.InboundTLSOptionsContainer{
				TLS: &tls,
			},
		},
	}
}
