package singboxserver

import (
	"fmt"
	"net/netip"
	"net/url"

	"github.com/awryme/unchained/app/appconfig"
	"github.com/awryme/unchained/pkg/protocols"
	"github.com/awryme/unchained/pkg/protocols/vless/vlessvision"
)

func MakeUrl(cfg appconfig.Unchained) (string, error) {
	switch cfg.Proto {
	case protocols.Trojan:
		return makeUrlTrojan(cfg), nil
	case protocols.Vless:
		return makeUrlVless(cfg), nil
	}

	return "", protocols.ErrInvalid(cfg.Proto)
}

func makeUrlTrojan(cfg appconfig.Unchained) string {
	u := &url.URL{
		Scheme:   protocols.Trojan,
		Host:     getHost(cfg.AppInfo, cfg.Listen),
		User:     url.User(cfg.TrojanPassword),
		Fragment: cfg.Name(),
	}
	q := u.Query()
	setRealityParams(q, cfg.Singbox.Reality)
	u.RawQuery = q.Encode()

	return u.String()
}

func makeUrlVless(cfg appconfig.Unchained) string {
	u := &url.URL{
		Scheme:   protocols.Vless,
		Host:     getHost(cfg.AppInfo, cfg.Listen),
		User:     url.User(cfg.VlessUUID.String()),
		Fragment: cfg.Name(),
	}
	q := u.Query()
	setRealityParams(q, cfg.Singbox.Reality)
	q.Set("flow", vlessvision.Flow)
	u.RawQuery = q.Encode()

	return u.String()
}

func getHost(appInfo appconfig.AppInfo, listen netip.AddrPort) string {
	return fmt.Sprintf("%s:%d", appInfo.PublicIP.String(), listen.Port())
}

func setRealityParams(q url.Values, reality appconfig.Reality) {
	q.Set("type", "tcp")
	q.Set("security", "reality")
	q.Set("fp", "chrome")
	q.Set("sni", reality.Server)
	q.Set("pbk", reality.PublicKey)
	q.Set("sid", reality.ShortId)
}
