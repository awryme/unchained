package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/trojan"
)

type CmdRun struct {
	LogLevel string `help:"sing-box log level" default:"${log_level}"`
	DNS      string `help:"sing-box dns" default:"${dns}"`
	Name     string `help:"sing-box proxy name (used to identify proxy in clients), random by default"`
	NoConfig bool   `help:"do not generate config and ignore existing"`
}

func (c *CmdRun) Run(app *App) error {
	ctx := context.Background()
	cfg, err := c.readOrGenerateConfig(app.Config)
	if err != nil {
		return err
	}
	instance, err := trojan.Run(ctx, cfg)
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
	fmt.Println("Got ctrl+c / interrupt, quitting")
	return instance.Close()
}

func (c *CmdRun) readOrGenerateConfig(file string) (cfg appconfig.Config, err error) {
	if c.NoConfig {
		return c.generateConfig()
	}
	rc := &appconfig.RuntimeConfig{
		LogLevel: c.LogLevel,
		DNS:      c.DNS,
		Name:     c.Name,
	}
	cfg, err = appconfig.Read(file, rc)
	if err == nil {
		return
	}
	if !errors.Is(err, os.ErrNotExist) {
		return
	}

	cfg, err = appconfig.Generate(rc)
	if err != nil {
		return
	}

	err = appconfig.Write(cfg, file)
	return
}

func (c *CmdRun) generateConfig() (cfg appconfig.Config, err error) {
	rc := &appconfig.RuntimeConfig{
		LogLevel: c.LogLevel,
		DNS:      c.DNS,
		Name:     c.Name,
	}

	return appconfig.Generate(rc)
}
