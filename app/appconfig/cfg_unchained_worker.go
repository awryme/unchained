package appconfig

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/awryme/unchained/pkg/xnet"
)

type UnchainedWorkerRuntimeParams struct {
	LogLevel string
	DNS      string
	Tags     []string
}

type UnchainedWorker struct {
	Singbox Singbox `json:"singbox"`
	AppInfo AppInfo `json:"app_info"`

	Worker Worker `json:"worker"`

	ListenVless  netip.AddrPort `json:"listen_vless"`
	ListenTrojan netip.AddrPort `json:"listen_trojan"`
}

// func (cfg UnchainedWorker) Name() string {
// 	var tags string
// 	if len(cfg.AppInfo.Tags) > 0 {
// 		tags = "_" + strings.Join(cfg.AppInfo.Tags, "_")
// 	}
// 	return fmt.Sprintf("%s%s_%s", cfg.AppInfo.ID, tags, cfg.Proto)
// }

func (cfg *UnchainedWorker) setRuntimeParams(params *UnchainedWorkerRuntimeParams) {
	if params == nil {
		return
	}

	trySet := func(value *string, param string) {
		if param != "" {
			*value = param
		}
	}

	trySet(&cfg.Singbox.LogLevel, params.LogLevel)
	trySet(&cfg.Singbox.DNS, params.DNS)

	if len(params.Tags) > 0 {
		cfg.AppInfo.Tags = params.Tags
	}
}

func (cfg *UnchainedWorker) Generate(ctx context.Context, params *UnchainedWorkerRuntimeParams) (err error) {
	cfg.Singbox.LogLevel = DefaultLogLevel
	cfg.Singbox.DNS = DefaultDns

	cfg.setRuntimeParams(params)

	cfg.ListenVless, err = xnet.GetRandomListenAddr(DefaultListenAddr)
	if err != nil {
		return fmt.Errorf("get random vless listen addr: %w", err)
	}

	cfg.ListenTrojan, err = xnet.GetRandomListenAddr(DefaultListenAddr)
	if err != nil {
		return fmt.Errorf("get random trojan listen addr: %w", err)
	}

	if err := cfg.AppInfo.generate(ctx); err != nil {
		return fmt.Errorf("generate app info: %w", err)
	}
	if err := cfg.Singbox.generate(ctx); err != nil {
		return fmt.Errorf("generate singbox config: %w", err)
	}

	if err := cfg.Worker.generate(cfg.AppInfo.PublicIP); err != nil {
		return fmt.Errorf("generate worker config: %w", err)
	}

	return nil
}
