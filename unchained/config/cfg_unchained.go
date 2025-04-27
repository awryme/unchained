package config

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/sethvargo/go-password/password"
)

type DynamicParams struct {
	LogLevel string
	DNS      string
	Proto    string
	Tags     []string
}

func (params *DynamicParams) Set(cfg *Unchained) {
	if params == nil || cfg == nil {
		return
	}
	cfg.Singbox.LogLevel = params.LogLevel
	cfg.Singbox.DNS = params.DNS
	cfg.Proto = params.Proto
	cfg.AppInfo.Tags = params.Tags
}

type Unchained struct {
	Singbox Singbox `json:"singbox"`
	AppInfo AppInfo `json:"app_info"`

	Proto          string    `json:"proto"`
	TrojanPassword string    `json:"trojan_password"`
	VlessUUID      uuid.UUID `json:"vless_uuid"`
}

func (cfg *Unchained) Generate(ctx context.Context, params *DynamicParams) (err error) {
	params.Set(cfg)

	if err := cfg.AppInfo.Generate(ctx); err != nil {
		return fmt.Errorf("generate app info: %w", err)
	}
	if err := cfg.Singbox.Generate(ctx); err != nil {
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
