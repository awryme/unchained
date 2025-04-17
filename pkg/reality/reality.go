package reality

import (
	"encoding/base64"
	"fmt"
	"math/rand/v2"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const DefaultServer = "github.com"
const DefaultServerPort = 443
const DefaultTimeDiff = time.Minute

func MakeRealityKeyPair() (privateKey string, publicKey string, err error) {
	wgPrivateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", "", err
	}

	wgPublicKey := wgPrivateKey.PublicKey()
	encode := func(key wgtypes.Key) string {
		return base64.RawURLEncoding.EncodeToString(key[:])
	}

	privateKey = encode(wgPrivateKey)
	publicKey = encode(wgPublicKey)

	return
}

func MakeShortId() string {
	sid := fmt.Sprintf("%x", rand.Uint32())
	if len(sid)%2 == 1 {
		sid += "f"
	}
	return sid
}
