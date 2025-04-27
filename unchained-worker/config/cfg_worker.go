package config

import (
	"fmt"
	"net/netip"

	"github.com/awryme/unchained/pkg/xnet"
	"github.com/awryme/unchained/unchained/defaults"
)

type Worker struct {
	Listen netip.AddrPort `json:"listen"`
}

func (w *Worker) Generate(publicIP netip.Addr) (err error) {
	w.Listen, err = xnet.GetRandomListenAddr(defaults.ListenAddr)
	if err != nil {
		return fmt.Errorf("get random vless listen addr: %w", err)
	}

	return nil
}
