package trojan

import (
	"fmt"
	"net/url"

	"github.com/awryme/unchained/appconfig"
)

func MakeUrl(cfg appconfig.Config) string {
	u := &url.URL{
		Scheme:   "trojan",
		Host:     fmt.Sprintf("%s:%d", cfg.PublicIp, cfg.Port),
		User:     url.User(cfg.Password),
		Fragment: cfg.Name,
	}
	q := u.Query()
	q.Set("type", "tcp")
	q.Set("security", "reality")
	q.Set("fp", "chrome")
	// q.Set("sni", "") // fixme: use sni?
	q.Set("pbk", cfg.RealityPublicKey)
	q.Set("sid", cfg.RealityShortId)
	u.RawQuery = q.Encode()

	return u.String()
}
