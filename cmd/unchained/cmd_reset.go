package main

import "os"

type CmdReset struct{}

func (c *CmdReset) Run(app *App) error {
	return os.Remove(app.Config)
}
