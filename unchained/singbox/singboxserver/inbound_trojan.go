package singboxserver

import (
	"net/netip"

	"github.com/awryme/unchained/pkg/protocols/trojan/trojaninbound"
	"github.com/awryme/unchained/pkg/protocols/trojan/trojanproto"
	"github.com/awryme/unchained/unchained/config"
	"github.com/sagernet/sing-box/option"
)

type InboundTrojan struct {
	listen    netip.AddrPort
	reality   config.Reality
	userStore trojanproto.UserStore
}

func NewInboundTrojan(params config.ProxyParams, userStore trojanproto.UserStore) InboundTrojan {
	return InboundTrojan{
		listen:    params.Listen,
		reality:   params.Reality,
		userStore: userStore,
	}
}
func (inb InboundTrojan) MakeInbound() option.Inbound {
	return trojaninbound.MakeInbound(
		inb.userStore,
		makeListenOptions(inb.listen),
		makeTlsOptions(inb.reality),
	)
}
