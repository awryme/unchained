package singboxserver

import (
	"fmt"
	"net/url"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/protocols"
	"github.com/awryme/unchained/pkg/protocols/vless/vlessvision"
)

func MakeUrl(cfg appconfig.Config) (string, error) {
	switch cfg.Proto {
	case protocols.Trojan:
		return makeUrlTrojan(cfg), nil
	case protocols.Vless:
		return makeUrlVless(cfg), nil
	}

	return "", protocols.ErrInvalid(cfg.Proto)
}

func makeUrlTrojan(cfg appconfig.Config) string {
	u := &url.URL{
		Scheme:   protocols.Trojan,
		Host:     getHost(cfg),
		User:     url.User(cfg.TrojanPassword),
		Fragment: cfg.Name(),
	}
	q := getCommonQuery(cfg)
	u.RawQuery = q.Encode()

	return u.String()
}

func makeUrlVless(cfg appconfig.Config) string {
	u := &url.URL{
		Scheme:   protocols.Vless,
		Host:     getHost(cfg),
		User:     url.User(cfg.VlessUUID.String()),
		Fragment: cfg.Name(),
	}
	u.Query()
	q := getCommonQuery(cfg)
	q.Set("flow", vlessvision.Flow)
	u.RawQuery = q.Encode()

	return u.String()
}

func getHost(cfg appconfig.Config) string {
	return fmt.Sprintf("%s:%d", cfg.PublicIP.String(), cfg.Listen.Port())
}

func getCommonQuery(cfg appconfig.Config) url.Values {
	q := make(url.Values)

	q.Set("type", "tcp")
	q.Set("security", "reality")
	q.Set("fp", "chrome")
	q.Set("sni", cfg.Reality.Server)
	q.Set("pbk", cfg.Reality.PublicKey)
	q.Set("sid", cfg.Reality.ShortId)

	return q
}
