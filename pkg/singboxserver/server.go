package singboxserver

import (
	"context"
	"fmt"
	"io"
	"net/netip"

	"github.com/awryme/unchained/appconfig"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
)

func Run(ctx context.Context, cfg appconfig.Config) (io.Closer, error) {
	addr, err := netip.ParseAddr(cfg.Listen.Addr)
	if err != nil {
		return nil, fmt.Errorf("parse ip addr to listen on: %w", err)
	}
	ctx = box.Context(ctx,
		include.InboundRegistry(),
		include.OutboundRegistry(),
		include.EndpointRegistry(),
	)

	inbound, err := makeInbountOptions(cfg, addr)
	if err != nil {
		return nil, fmt.Errorf("make proxy inbound: %w", err)
	}

	instance, err := box.New(box.Options{
		Context: ctx,
		Options: option.Options{
			Log:      makeLogOptions(cfg),
			DNS:      makeDnsOptions(cfg),
			Inbounds: []option.Inbound{inbound},
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
