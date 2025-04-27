package config

import (
	"context"
	"fmt"

	"github.com/awryme/unchained/unchained/config"
)

type DynamicParams struct {
	LogLevel string
	DNS      string
	Tags     []string
}

func (params *DynamicParams) Set(cfg *UnchainedWorker) {
	if params == nil || cfg == nil {
		return
	}
	cfg.Singbox.LogLevel = params.LogLevel
	cfg.Singbox.DNS = params.DNS
	cfg.AppInfo.Tags = params.Tags
}

type UnchainedWorker struct {
	Singbox config.Singbox `json:"singbox"`
	AppInfo config.AppInfo `json:"app_info"`

	Worker Worker `json:"worker"`
}

func (cfg *UnchainedWorker) Generate(ctx context.Context, params *DynamicParams) (err error) {
	params.Set(cfg)

	if err := cfg.AppInfo.Generate(ctx); err != nil {
		return fmt.Errorf("generate app info: %w", err)
	}
	if err := cfg.Singbox.Generate(ctx); err != nil {
		return fmt.Errorf("generate singbox config: %w", err)
	}

	if err := cfg.Worker.Generate(cfg.AppInfo.PublicIP); err != nil {
		return fmt.Errorf("generate worker config: %w", err)
	}

	return nil
}
