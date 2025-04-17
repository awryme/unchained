package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/awryme/unchained/app/appconfig"
	"github.com/awryme/unchained/app/clilog"
	"github.com/awryme/unchained/app/singbox/memoryuserstore"
	"github.com/awryme/unchained/app/singbox/singboxserver"
	"github.com/awryme/unchained/pkg/protocols"
)

type CmdRun struct {
	RuntimeParams `embed:""`

	NoConfig bool `help:"only generate config, ignore existing"`
}

func (c *CmdRun) Run(app *App) error {
	ctx := context.Background()
	cfg, err := c.getConfig(ctx, app.Config)
	if err != nil {
		return err
	}

	inbound, err := c.getInbound(cfg)
	if err != nil {
		return err
	}

	instance, err := singboxserver.Run(ctx, cfg.Singbox, inbound)
	if err != nil {
		return err
	}
	clilog.Log("Started at", time.Now().Format(time.DateTime))
	err = printInfo(cfg)
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	clilog.Log("Got ctrl+c / interrupt, quitting")
	return instance.Close()
}

func (c *CmdRun) getInbound(cfg appconfig.Unchained) (singboxserver.InboundMaker, error) {
	switch cfg.Proto {
	case protocols.Trojan:
		userStore := memoryuserstore.NewTrojan()
		userStore.Add("Single user", cfg.TrojanPassword)
		return singboxserver.NewInboundTrojan(cfg.Listen, cfg.Singbox, userStore), nil
	case protocols.Vless:
		userStore := memoryuserstore.NewVless()
		userStore.Add("Single user", cfg.VlessUUID)
		return singboxserver.NewInboundVless(cfg.Listen, cfg.Singbox, userStore), nil
	}

	return nil, protocols.ErrInvalid(cfg.Proto)
}

func (c *CmdRun) getConfig(ctx context.Context, file string) (cfg appconfig.Unchained, err error) {
	params := c.GetRuntimeParams()
	if c.NoConfig {
		err := cfg.Generate(ctx, params)
		return cfg, err
	}

	cfg, err = appconfig.ReadUnchained(file, params)
	if errors.Is(err, os.ErrNotExist) {
		// cfg file not found, generate new one
		err = cfg.Generate(ctx, params)
	}
	if err != nil {
		return cfg, err
	}

	// write any changes that we applied
	err = appconfig.WriteUnchained(cfg, file)
	return cfg, err
}
