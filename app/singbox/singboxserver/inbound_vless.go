package singboxserver

import (
	"net/netip"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/protocols/vless/vlessinbound"
	"github.com/awryme/unchained/pkg/protocols/vless/vlessproto"
	"github.com/sagernet/sing-box/option"
)

type InboundVless struct {
	listen    netip.AddrPort
	reality   appconfig.Reality
	userStore vlessproto.UserStore
}

func NewInboundVless(listen netip.AddrPort, reality appconfig.Reality, userStore vlessproto.UserStore) InboundVless {
	return InboundVless{
		listen:    listen,
		reality:   reality,
		userStore: userStore,
	}
}

func (inb InboundVless) MakeInbound() option.Inbound {
	return vlessinbound.MakeInbound(
		inb.userStore,
		makeListenOptions(inb.listen),
		makeTlsOptions(inb.reality),
	)
}
