package config

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/awryme/unchained/pkg/xnet"
	"github.com/awryme/unchained/unchained/defaults"
)

type ProxyParams struct {
	Listen  netip.AddrPort `json:"listen"`
	Reality Reality        `json:"reality"`
}

func (cfg *ProxyParams) Generate(ctx context.Context) (err error) {
	cfg.Listen, err = xnet.GetRandomListenAddr(defaults.ListenAddr)
	if err != nil {
		return fmt.Errorf("get random listen addr: %w", err)
	}
	if err := cfg.Reality.Generate(); err != nil {
		return err
	}

	return nil
}
