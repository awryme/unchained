package urlmaker

import (
	"fmt"
	"net/url"

	"github.com/awryme/unchained/pkg/protocols/vless/vlessvision"
	"github.com/awryme/unchained/unchained/config"
	"github.com/awryme/unchained/unchained/protocols"
	"github.com/gofrs/uuid/v5"
)

type UrlMaker struct {
	appInfo config.AppInfo
	singbox config.Singbox
}

func New(appInfo config.AppInfo, singbox config.Singbox) *UrlMaker {
	return &UrlMaker{appInfo, singbox}
}

func (m *UrlMaker) MakeByProto(proto string, vlesID uuid.UUID, trojanPassword string) (string, error) {
	switch proto {
	case protocols.Trojan:
		return m.MakeTrojan(trojanPassword), nil
	case protocols.Vless:
		return m.MakeVless(vlesID), nil
	}

	return "", protocols.ErrInvalid(proto)

}

func (m *UrlMaker) MakeTrojan(password string) string {
	appInfo := m.appInfo
	proxyParams := m.singbox.TrojanProxy
	u := &url.URL{
		Scheme:   protocols.Trojan,
		Host:     fmt.Sprintf("%s:%d", appInfo.PublicIP.String(), proxyParams.Listen.Port()),
		User:     url.User(password),
		Fragment: appInfo.Name(protocols.Trojan),
	}
	q := u.Query()
	setRealityParams(q, proxyParams.Reality)
	u.RawQuery = q.Encode()

	return u.String()
}

func (m *UrlMaker) MakeVless(vlessID uuid.UUID) string {
	appInfo := m.appInfo
	proxyParams := m.singbox.VlessProxy

	u := &url.URL{
		Scheme:   protocols.Vless,
		Host:     fmt.Sprintf("%s:%d", appInfo.PublicIP.String(), proxyParams.Listen.Port()),
		User:     url.User(vlessID.String()),
		Fragment: appInfo.Name(protocols.Vless),
	}
	q := u.Query()
	setRealityParams(q, proxyParams.Reality)
	q.Set("flow", vlessvision.Flow)
	u.RawQuery = q.Encode()

	return u.String()

}

func setRealityParams(q url.Values, reality config.Reality) {
	q.Set("type", "tcp")
	q.Set("security", "reality")
	q.Set("fp", "chrome")
	q.Set("sni", reality.Server)
	q.Set("pbk", reality.PublicKey)
	q.Set("sid", reality.ShortId)
}
