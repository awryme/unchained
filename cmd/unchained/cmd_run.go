package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/awryme/unchained/unchained/clilog"
	"github.com/awryme/unchained/unchained/config"
	"github.com/awryme/unchained/unchained/protocols"
	"github.com/awryme/unchained/unchained/singbox/memoryuserstore"
	"github.com/awryme/unchained/unchained/singbox/singboxserver"
)

type CmdRun struct {
	LogLevel string   `help:"sing-box log level" default:"${log_level}"`
	DNS      string   `help:"dns address, in sing-box format" default:"${dns}"`
	Proto    string   `short:"p" help:"set used protocol: ${enum}" enum:"${protos}" default:"${default_proto}"`
	Tags     []string `help:"proxy tags (used to identify proxy in client apps)"`

	NoConfig bool   `help:"only generate config, ignore existing"`
	Dir      string `help:"dir to store config file and data" default:"./data/"`
}

func (cmd *CmdRun) Run(app *App) error {
	ctx := context.Background()

	cfg, err := cmd.getConfig(ctx)
	if err != nil {
		return err
	}

	inbound, err := cmd.getInbound(cfg)
	if err != nil {
		return err
	}

	instance, err := singboxserver.Run(ctx, cfg.Singbox, inbound)
	if err != nil {
		return err
	}
	clilog.Log("Started at", time.Now().Format(time.DateTime))
	err = printInfo(cmd.Dir, cfg, cfg.AppInfo, cfg.Singbox)
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	clilog.Log("Got ctrl+c / interrupt, quitting")
	return instance.Close()
}

func (cmd *CmdRun) getInbound(cfg config.Unchained) (singboxserver.InboundMaker, error) {
	switch cfg.Proto {
	case protocols.Trojan:
		userStore := memoryuserstore.NewTrojan()
		userStore.Add("Single user", cfg.TrojanPassword)
		return singboxserver.NewInboundTrojan(cfg.Singbox.TrojanProxy, userStore), nil
	case protocols.Vless:
		userStore := memoryuserstore.NewVless()
		userStore.Add("Single user", cfg.VlessUUID)
		return singboxserver.NewInboundVless(cfg.Singbox.VlessProxy, userStore), nil
	}

	return nil, protocols.ErrInvalid(cfg.Proto)
}

func (cmd *CmdRun) getConfigName() (string, error) {
	if err := os.MkdirAll(cmd.Dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("make data dir: %w", err)
	}
	file := filepath.Join(cmd.Dir, ConfigName)
	return file, nil
}

func (cmd *CmdRun) getConfig(ctx context.Context) (cfg config.Unchained, err error) {
	params := &config.DynamicParams{
		LogLevel: cmd.LogLevel,
		DNS:      cmd.DNS,
		Proto:    cmd.Proto,
		Tags:     cmd.Tags,
	}
	if cmd.NoConfig {
		err := cfg.Generate(ctx, params)
		return cfg, err
	}

	file, err := cmd.getConfigName()
	if err != nil {
		return cfg, fmt.Errorf("get config: %w", err)
	}
	cfg, err = config.Read(file, params)
	if errors.Is(err, os.ErrNotExist) {
		// cfg file not found, generate new one
		err = cfg.Generate(ctx, params)
	}
	if err != nil {
		return cfg, err
	}

	// write any changes that we applied
	err = config.Write(cfg, file)
	return cfg, err
}
