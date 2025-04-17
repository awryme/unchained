package appconfig

import (
	"context"
	"fmt"
	"net/netip"
	"strings"

	"github.com/awryme/unchained/pkg/protocols"
	"github.com/awryme/unchained/pkg/xnet"
	"github.com/gofrs/uuid/v5"
	"github.com/sethvargo/go-password/password"
)

type UnchainedRuntimeParams struct {
	LogLevel string
	DNS      string
	Proto    string
	Tags     []string
}

type Unchained struct {
	Singbox Singbox `json:"singbox"`
	AppInfo AppInfo `json:"app_info"`

	Listen netip.AddrPort `json:"listen"`

	Proto          string    `json:"proto"`
	TrojanPassword string    `json:"trojan_password"`
	VlessUUID      uuid.UUID `json:"vless_uuid"`
}

func (cfg Unchained) Name() string {
	var tags string
	if len(cfg.AppInfo.Tags) > 0 {
		tags = "_" + strings.Join(cfg.AppInfo.Tags, "_")
	}
	return fmt.Sprintf("%s%s_%s", cfg.AppInfo.ID, tags, cfg.Proto)
}

func (cfg *Unchained) setRuntimeParams(params *UnchainedRuntimeParams) {
	if params == nil {
		return
	}

	trySet := func(value *string, param string) {
		if param != "" {
			*value = param
		}
	}

	trySet(&cfg.Proto, params.Proto)
	trySet(&cfg.Singbox.LogLevel, params.LogLevel)
	trySet(&cfg.Singbox.DNS, params.DNS)

	if len(params.Tags) > 0 {
		cfg.AppInfo.Tags = params.Tags
	}
}

func (cfg *Unchained) Generate(ctx context.Context, params *UnchainedRuntimeParams) (err error) {
	cfg.Singbox.LogLevel = DefaultLogLevel
	cfg.Singbox.DNS = DefaultDns
	cfg.Proto = protocols.Vless

	cfg.setRuntimeParams(params)

	cfg.Listen, err = xnet.GetRandomListenAddr(DefaultListenAddr)
	if err != nil {
		return fmt.Errorf("get random listen addr: %w", err)
	}

	if err := cfg.AppInfo.generate(ctx); err != nil {
		return fmt.Errorf("generate app info: %w", err)
	}
	if err := cfg.Singbox.generate(ctx); err != nil {
		return fmt.Errorf("generate singbox config: %w", err)
	}

	if err := cfg.setTrojanPassword(); err != nil {
		return fmt.Errorf("generate trojan password: %w", err)
	}

	if err := cfg.setVlessUUID(); err != nil {
		return fmt.Errorf("generate vless uuid: %w", err)
	}

	return nil
}

func (cfg *Unchained) setTrojanPassword() error {
	const length = 16

	pwd, err := password.Generate(length, length/3, 0, false, false)
	if err != nil {
		return err
	}

	cfg.TrojanPassword = pwd
	return nil
}

func (cfg *Unchained) setVlessUUID() error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	cfg.VlessUUID = id
	return nil
}
