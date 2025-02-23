package appconfig

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"time"

	"github.com/awryme/unchained/constants"
	"github.com/awryme/unchained/pkg/clilog"
	"github.com/awryme/unchained/pkg/ipv6detector"
	"github.com/awryme/unchained/pkg/protocols"
	"github.com/gofrs/uuid/v5"
	"github.com/sethvargo/go-password/password"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func Generate(ctx context.Context, params *RuntimeParams) (Config, error) {
	cfg := Config{
		LogLevel: constants.DefaultLogLevel,
		DNS:      constants.DefaultDns,
		Proto:    protocols.Trojan,
		Listen: Listen{
			Addr: constants.DefaultListenAddr,
		},
	}

	if err := cfg.setRandomID(); err != nil {
		return cfg, fmt.Errorf("generate random name: %w", err)
	}

	setRuntimeParams(&cfg, params)

	if err := cfg.setRandomPort(); err != nil {
		return cfg, fmt.Errorf("set random port: %w", err)
	}

	if err := cfg.setPublicIP(ctx); err != nil {
		return cfg, fmt.Errorf("set public ip: %w", err)
	}

	if err := cfg.setTrojanPassword(); err != nil {
		return cfg, fmt.Errorf("generate random password: %w", err)
	}

	if err := cfg.setVlessUUID(); err != nil {
		return cfg, fmt.Errorf("generate vless uuid: %w", err)
	}

	if err := cfg.setRealityConfig(); err != nil {
		return cfg, err
	}

	cfg.setIPv4Only(ctx)

	return cfg, nil
}

func (cfg *Config) setRealityConfig() error {
	privateKey, publicKey, err := generateRealityKeyPair()
	if err != nil {
		return fmt.Errorf("generate reality keypair: %w", err)
	}

	cfg.Reality = Reality{
		PrivateKey: privateKey,
		PublicKey:  publicKey,

		Server:   constants.DefaultRealityServer,
		ShortId:  generateRealityShortId(),
		TimeDiff: constants.DefaultRealityTimeDiff,
	}
	return nil
}

func (cfg *Config) setRandomPort() (err error) {
	addr := cfg.Listen.Addr

	a, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", addr))
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		return err
	}

	err = l.Close()
	if err != nil {
		return err
	}

	cfg.Listen.Port = l.Addr().(*net.TCPAddr).Port
	return nil
}

func (cfg *Config) setRandomID() error {
	const length = 6

	id, err := password.Generate(length, length/3, 0, true, false)
	if err != nil {
		return err
	}

	cfg.ID = id
	return nil
}

func (cfg *Config) setTrojanPassword() error {
	const length = 16

	pwd, err := password.Generate(length, length/3, 0, false, false)
	if err != nil {
		return err
	}

	cfg.TrojanPassword = pwd
	return nil
}

func (cfg *Config) setVlessUUID() error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	cfg.VlessUUID = id.String()
	return nil
}

func generateRealityKeyPair() (privateKey string, publicKey string, err error) {
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

func generateRealityShortId() string {
	sid := fmt.Sprintf("%x", rand.Uint32())
	if len(sid)%2 == 1 {
		sid += "f"
	}
	return sid
}

func (cfg *Config) setPublicIP(ctx context.Context) error {
	cl := &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, constants.IpifyUrl, nil)
	if err != nil {
		return err
	}

	res, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code %d", res.StatusCode)
	}

	ipBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	cfg.PublicIP = string(ipBytes)
	return nil
}

func (cfg *Config) setIPv4Only(ctx context.Context) {
	logf := func(format string, args ...any) {
		format = "detect ipv6: " + format
		clilog.Logf(format, args...)
	}
	hasIPv6 := ipv6detector.Detect(ctx, logf)
	if !hasIPv6 {
		clilog.Logf("ipv6 DNS is disabled")
		cfg.DNSIPv4Only = true
	}
}
