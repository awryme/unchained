package main

import (
	"context"
	"fmt"
	"os"

	"github.com/awryme/unchained/appconfig"
)

type CmdGenerate struct {
	RuntimeParams `embed:""`
}

func (c *CmdGenerate) Run(app *App) error {
	ctx := context.Background()
	if fileExists(app.Config) {
		return fmt.Errorf("config file %s already exists, run command 'reset' to cleanup, or remove manually", app.Config)
	}

	cfg, err := appconfig.Generate(ctx, c.GetRuntimeParams())
	if err != nil {
		return fmt.Errorf("generate config: %w", err)
	}

	return appconfig.Write(cfg, app.Config)
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
