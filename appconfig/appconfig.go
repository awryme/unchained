package appconfig

import (
	"fmt"
	"strings"
	"time"
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
	var tags string
	if len(cfg.Tags) > 0 {
		tags = "_" + strings.Join(cfg.Tags, "_")
	}
	return fmt.Sprintf("%s%s_%s", cfg.ID, tags, cfg.Proto)
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
