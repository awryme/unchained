package singboxserver

import (
	"net/netip"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/protocols/trojan/trojaninbound"
	"github.com/awryme/unchained/pkg/protocols/trojan/trojanproto"
	"github.com/sagernet/sing-box/option"
)

type InboundTrojan struct {
	listen    netip.AddrPort
	reality   appconfig.Reality
	userStore trojanproto.UserStore
}

func NewInboundTrojan(listen netip.AddrPort, reality appconfig.Reality, userStore trojanproto.UserStore) InboundTrojan {
	return InboundTrojan{
		listen:    listen,
		reality:   reality,
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
