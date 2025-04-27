package singboxserver

import (
	"net/netip"

	"github.com/awryme/unchained/pkg/protocols/vless/vlessinbound"
	"github.com/awryme/unchained/pkg/protocols/vless/vlessproto"
	"github.com/awryme/unchained/unchained/config"
	"github.com/sagernet/sing-box/option"
)

type InboundVless struct {
	listen    netip.AddrPort
	reality   config.Reality
	userStore vlessproto.UserStore
}

func NewInboundVless(params config.ProxyParams, userStore vlessproto.UserStore) InboundVless {
	return InboundVless{
		listen:    params.Listen,
		reality:   params.Reality,
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
