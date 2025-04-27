package singboxserver

import (
	"context"
	"fmt"
	"io"

	"github.com/awryme/unchained/pkg/protocols/trojan/trojaninbound"
	"github.com/awryme/unchained/pkg/protocols/vless/vlessinbound"
	"github.com/awryme/unchained/unchained/config"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/adapter/inbound"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
)

func makeInbountRegistry() *inbound.Registry {
	reg := include.InboundRegistry()

	vlessinbound.Register(reg)
	trojaninbound.Register(reg)

	return reg
}

func Run(ctx context.Context, cfg config.Singbox, inbounds ...InboundMaker) (io.Closer, error) {
	ctx = box.Context(ctx,
		makeInbountRegistry(),
		include.OutboundRegistry(),
		include.EndpointRegistry(),
	)

	singboxInbounds := make([]option.Inbound, len(inbounds))
	for i, maker := range inbounds {
		singboxInbounds[i] = maker.MakeInbound()
	}

	instance, err := box.New(box.Options{
		Context: ctx,
		Options: option.Options{
			Log:      makeLogOptions(cfg),
			DNS:      makeDnsOptions(cfg),
			Inbounds: singboxInbounds,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("make sing-box instance: %w", err)
	}
	err = instance.Start()
	if err != nil {
		return nil, fmt.Errorf("start sing-box instance: %w", err)
	}
	return instance, nil
}
