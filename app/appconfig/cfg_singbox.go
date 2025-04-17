package appconfig

import (
	"context"
	"fmt"

	"github.com/awryme/ipinfo"
	"github.com/awryme/unchained/app/clilog"
)

type Singbox struct {
	LogLevel    string `json:"log_level"`
	DNS         string `json:"dns"`
	DNSIPv4Only bool   `json:"dns_ipv4_only"`

	Reality Reality `json:"reality"`
}

func (cfg *Singbox) generate(ctx context.Context) (err error) {
	if err := cfg.setIpV4Only(ctx); err != nil {
		return fmt.Errorf("set public ip: %w", err)
	}

	if err := cfg.Reality.generate(); err != nil {
		return err
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
