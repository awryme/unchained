package singboxserver

import (
	"fmt"
	"net/url"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/protocols"
)

func MakeUrl(cfg appconfig.Config) (string, error) {
	switch cfg.Proto {
	case protocols.Trojan:
		return makeUrlTrojan(cfg), nil
	case protocols.Vless:
		return makeUrlVless(cfg), nil
	}

	return "", protocols.Invalid(cfg.Proto)
}

func makeUrlTrojan(cfg appconfig.Config) string {
	u := &url.URL{
		Scheme:   protocols.Trojan,
		Host:     fmt.Sprintf("%s:%d", cfg.PublicIP, cfg.Listen.Port),
		User:     url.User(cfg.TrojanPassword),
		Fragment: cfg.Name(),
	}
	q := u.Query()
	q.Set("type", "tcp")
	q.Set("security", "reality")
	q.Set("fp", "chrome")
	// q.Set("sni", "") // fixme: use sni?
	q.Set("pbk", cfg.Reality.PublicKey)
	q.Set("sid", cfg.Reality.ShortId)
	u.RawQuery = q.Encode()

	return u.String()
}

func makeUrlVless(cfg appconfig.Config) string {
	u := &url.URL{
		Scheme:   protocols.Vless,
		Host:     fmt.Sprintf("%s:%d", cfg.PublicIP, cfg.Listen.Port),
		User:     url.User(cfg.VlessUUID),
		Fragment: cfg.Name(),
	}
	q := u.Query()
	q.Set("flow", vlessFlowVision)
	q.Set("type", "tcp")
	q.Set("security", "reality")
	q.Set("fp", "chrome")
	// q.Set("sni", "") // fixme: use sni?
	q.Set("pbk", cfg.Reality.PublicKey)
	q.Set("sid", cfg.Reality.ShortId)
	u.RawQuery = q.Encode()

	return u.String()
}
