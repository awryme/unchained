package vlessinbound

import (
	"github.com/awryme/unchained/pkg/protocols/vless/vlessproto"
	"github.com/sagernet/sing-box/option"
)

type VLESSInboundOptions struct {
	option.ListenOptions
	UserStore vlessproto.UserStore `json:"users,omitempty"`
	option.InboundTLSOptionsContainer
	Multiplex *option.InboundMultiplexOptions `json:"multiplex,omitempty"`
	Transport *option.V2RayTransportOptions   `json:"transport,omitempty"`
}

func MakeInbound(userStore vlessproto.UserStore, listen option.ListenOptions, tls option.InboundTLSOptions) option.Inbound {
	return option.Inbound{
		Type: TypeVlessV2,
		Tag:  "inbound-" + TypeVlessV2,
		Options: &VLESSInboundOptions{
			ListenOptions: listen,
			UserStore:     userStore,
			InboundTLSOptionsContainer: option.InboundTLSOptionsContainer{
				TLS: &tls,
			},
		},
	}
}
