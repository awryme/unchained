package config

import (
	"context"
	"fmt"

	"github.com/awryme/ipinfo"
	"github.com/awryme/unchained/unchained/clilog"
)

type Singbox struct {
	LogLevel    string `json:"log_level"`
	DNS         string `json:"dns"`
	DNSIPv4Only bool   `json:"dns_ipv4_only"`

	VlessProxy  ProxyParams `json:"vless_proxy"`
	TrojanProxy ProxyParams `json:"trojan_proxy"`
}

func (cfg *Singbox) Generate(ctx context.Context) (err error) {
	if err := cfg.setIpV4Only(ctx); err != nil {
		return fmt.Errorf("set public ip: %w", err)
	}

	if err := cfg.VlessProxy.Generate(ctx); err != nil {
		return fmt.Errorf("generate vless proxy: %w", err)
	}

	if err := cfg.TrojanProxy.Generate(ctx); err != nil {
		return fmt.Errorf("generate trojan proxy: %w", err)
	}

	return nil
}

func (cfg *Singbox) setIpV4Only(ctx context.Context) error {
	// detect ipv6, set DNSIPv4Only
	// no errors, just log
	_, err := ipinfo.PublicIPv6(ctx)
	if err != nil {
		cfg.DNSIPv4Only = true
		clilog.Log("ipv6 disabled, err:", err)
		return nil
	}
	cfg.DNSIPv4Only = false
	return nil
}
