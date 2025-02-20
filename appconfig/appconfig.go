package appconfig

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/awryme/unchained/constants"
	"github.com/sethvargo/go-password/password"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Config struct {
	LogLevel          string        `json:"log_level"`
	DNS               string        `json:"dns"`
	ListenAddr        string        `json:"addr"`
	Port              int           `json:"port"`
	PublicIp          string        `json:"public_ip"`
	Name              string        `json:"name"`
	Password          string        `json:"password"`
	RealityServer     string        `json:"reality_server"`
	RealityPrivateKey string        `json:"reality_private_key"`
	RealityPublicKey  string        `json:"reality_public_key"`
	RealityShortId    string        `json:"reality_short_id"`
	RealityTimeDiff   time.Duration `json:"reality_time_diff"`
}

type RuntimeConfig struct {
	LogLevel string
	DNS      string
	Name     string
}

func Read(file string, rc *RuntimeConfig) (Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return Config{}, fmt.Errorf("open config file: %w", err)
	}
	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("decode config file: %w", err)
	}
	setRuntimeCfg(&cfg, rc)
	return cfg, nil
}

func Write(cfg Config, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("create config file: %w", err)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(cfg)
	if err != nil {
		return fmt.Errorf("encode config file: %w", err)
	}
	return nil
}

func setRuntimeCfg(cfg *Config, rc *RuntimeConfig) {
	if rc == nil {
		return
	}
	if rc.LogLevel != "" {
		cfg.LogLevel = rc.LogLevel
	}
	if rc.DNS != "" {
		cfg.DNS = rc.DNS
	}
	if rc.Name != "" {
		cfg.Name = rc.Name
	}
}

func Generate(rc *RuntimeConfig) (Config, error) {
	port, err := getRandomPort(constants.ListenAddr)
	if err != nil {
		return Config{}, fmt.Errorf("get random port: %w", err)
	}
	name, err := generateRandomName()
	if err != nil {
		return Config{}, fmt.Errorf("generate random name: %w", err)
	}
	password, err := generateRandomPassword()
	if err != nil {
		return Config{}, fmt.Errorf("generate random password: %w", err)
	}

	privateKey, publicKey, err := generateRealityKeyPair()
	if err != nil {
		return Config{}, fmt.Errorf("generate reality keypair: %w", err)
	}

	shortId := generateRealityShortId()

	publicIp, err := getPublicIp()
	if err != nil {
		return Config{}, fmt.Errorf("get public ip: %w", err)
	}

	cfg := Config{
		LogLevel:          constants.DefaultLogLevel,
		DNS:               constants.DefaultDns,
		ListenAddr:        constants.ListenAddr,
		Port:              port,
		Name:              name,
		Password:          password,
		PublicIp:          publicIp,
		RealityServer:     constants.DefaultRealityServer,
		RealityPrivateKey: privateKey,
		RealityPublicKey:  publicKey,
		RealityShortId:    shortId,
		RealityTimeDiff:   constants.DefaultRealityTimeDiff,
	}
	setRuntimeCfg(&cfg, rc)
	return cfg, nil
}

func getRandomPort(addr string) (port int, err error) {
	a, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", addr))
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		return 0, err
	}
	err = l.Close()
	if err != nil {
		return 0, err
	}
	return l.Addr().(*net.TCPAddr).Port, nil
}

func generateRandomName() (string, error) {
	const length = 6
	return password.Generate(length, length/3, 0, true, false)
}

func generateRandomPassword() (string, error) {
	const length = 12
	return password.Generate(length, length/3, 0, false, false)
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
	return fmt.Sprintf("%x", rand.Uint32())
}

func getPublicIp() (string, error) {
	cl := &http.Client{
		Timeout: time.Second * 30,
	}
	res, err := cl.Get(constants.IpifyUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code %d", res.StatusCode)
	}

	ipBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(ipBytes), nil
}
