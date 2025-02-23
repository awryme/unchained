package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/clilog"
	"github.com/awryme/unchained/pkg/singboxserver"
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
	instance, err := singboxserver.Run(ctx, cfg)
	if err != nil {
		return err
	}
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

func (c *CmdRun) getConfig(ctx context.Context, file string) (appconfig.Config, error) {
	params := c.GetRuntimeParams()
	if c.NoConfig {
		return appconfig.Generate(ctx, params)
	}

	cfg, err := appconfig.Read(file, params)
	if errors.Is(err, os.ErrNotExist) {
		cfg, err = appconfig.Generate(ctx, params)
	}
	if err != nil {
		return cfg, err
	}

	err = appconfig.Write(cfg, file)
	return cfg, err
}
