package appconfig

import (
	"fmt"
	"strings"
	"time"

	"github.com/awryme/unchained/constants"
)

type Config struct {
	LogLevel       string   `json:"log_level"`
	DNS            string   `json:"dns"`
	DNSIPv4Only    bool     `json:"dns_ipv4_only"`
	Proto          string   `json:"proto"`
	ID             string   `json:"id"`
	Tags           []string `json:"tags"`
	PublicIP       string   `json:"public_ip"`
	TrojanPassword string   `json:"trojan_password"`
	VlessUUID      string   `json:"vless_uuid"`

	Listen  Listen  `json:"listen"`
	Reality Reality `json:"reality"`
}

func (cfg Config) Name() string {
	var tagSuffix string
	if len(cfg.Tags) > 0 {
		tagSuffix = "_" + strings.Join(cfg.Tags, "_")
	}
	return fmt.Sprintf("%s-%s%s", constants.AppName, cfg.ID, tagSuffix)
}

type Listen struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

type Reality struct {
	Server     string        `json:"server"`
	PrivateKey string        `json:"private_key"`
	PublicKey  string        `json:"public_key"`
	ShortId    string        `json:"short_id"`
	TimeDiff   time.Duration `json:"time_diff"`
}
