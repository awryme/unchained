package main

import (
	"fmt"

	"github.com/awryme/unchained/appconfig"
)

type CmdPrint struct {
}

func (c *CmdPrint) Run(app *App) error {
	cfg, err := appconfig.Read(app.Config, nil)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	return printInfo(cfg)
}
