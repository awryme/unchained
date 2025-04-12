package singboxserver

import "github.com/sagernet/sing-box/option"

type InboundMaker interface {
	MakeInbound() option.Inbound
}
