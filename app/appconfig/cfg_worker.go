package appconfig

import (
	"encoding/base64"
	"fmt"
	"net/netip"

	"github.com/awryme/unchained/pkg/certgen"
	"github.com/awryme/unchained/pkg/xnet"
	"github.com/gofrs/uuid/v5"
	"github.com/sethvargo/go-password/password"
)

type Worker struct {
	ID     uuid.UUID      `json:"worker_id"`
	Listen netip.AddrPort `json:"listen_worker"`

	JwtSecret string `json:"jwt_secret"`

	EncodedCert    string `json:"cert"`
	EncodedCertKey string `json:"cert_key"`
}

func (w *Worker) generate(publicIP netip.Addr) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate worker id: %w", err)
	}
	w.ID = id

	w.Listen, err = xnet.GetRandomListenAddr(DefaultListenAddr)
	if err != nil {
		return fmt.Errorf("get random vless listen addr: %w", err)
	}

	const length = 32
	secret, err := password.Generate(length, length/3, length/3, false, false)
	if err != nil {
		return err
	}
	w.JwtSecret = secret

	crtpem, keypem, err := certgen.Generate(publicIP)
	if err != nil {
		return fmt.Errorf("generate certificate for api server: %w", err)
	}

	w.EncodedCert = base64.URLEncoding.EncodeToString(crtpem)
	w.EncodedCertKey = base64.URLEncoding.EncodeToString(keypem)

	return nil
}
