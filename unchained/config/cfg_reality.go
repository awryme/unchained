package config

import (
	"fmt"
	"time"

	"github.com/awryme/unchained/pkg/reality"
)

type Reality struct {
	Server     string        `json:"server"`
	PrivateKey string        `json:"private_key"`
	PublicKey  string        `json:"public_key"`
	ShortId    string        `json:"short_id"`
	TimeDiff   time.Duration `json:"time_diff"`
}

func (r *Reality) Generate() error {
	privateKey, publicKey, err := reality.MakeRealityKeyPair()
	if err != nil {
		return fmt.Errorf("generate reality keypair: %w", err)
	}
	r.PrivateKey = privateKey
	r.PublicKey = publicKey

	r.ShortId = reality.MakeShortId()

	r.Server = reality.DefaultServer
	r.TimeDiff = reality.DefaultTimeDiff

	return nil
}
