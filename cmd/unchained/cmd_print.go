package main

import (
	"fmt"
	"path/filepath"

	"github.com/awryme/unchained/unchained/config"
)

type CmdPrint struct {
	Dir string `help:"dir to store config file and data" default:"./data/"`
}

func (cmd *CmdPrint) Run(app *App) error {
	file := filepath.Join(cmd.Dir, ConfigName)
	cfg, err := config.Read(file, nil)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	return printInfo(cmd.Dir, cfg, cfg.AppInfo, cfg.Singbox)
}
