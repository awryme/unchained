package main

import (
	"fmt"

	"github.com/awryme/unchained/unchained/config"
)

type CmdPrint struct {
}

func (c *CmdPrint) Run(app *App) error {
	cfg, err := config.Read(app.Config, nil)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	return printInfo(cfg, cfg.AppInfo, cfg.Singbox)
}
