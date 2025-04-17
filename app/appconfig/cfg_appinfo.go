package appconfig

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/awryme/ipinfo"
	"github.com/sethvargo/go-password/password"
)

type AppInfo struct {
	ID       string     `json:"id"`
	Tags     []string   `json:"tags"`
	PublicIP netip.Addr `json:"public_ip"`
}

func (cfg *AppInfo) generate(ctx context.Context) error {
	// set ipv4
	ip, err := ipinfo.PublicIPv4(ctx)
	if err != nil {
		return fmt.Errorf("get public ip: %w", err)
	}
	const length = 6

	// generate id
	id, err := password.Generate(length, length/3, 0, true, false)
	if err != nil {
		return err
	}

	cfg.ID = id
	cfg.PublicIP = ip

	return nil
}
