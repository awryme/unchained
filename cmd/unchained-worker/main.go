package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/awryme/unchained/unchained/defaults"
	"github.com/awryme/unchained/unchained/protocols"
)

const AppName = "unchained-worker"
const ConfigName = AppName + ".json"

type App struct {
	Run CmdRun `cmd:"" help:"run vpn server worker, generate config if it doesn't exist"`
}

func main() {
	var app App
	appctx := kong.Parse(&app,
		kong.Name(AppName),
		kong.Description(fmt.Sprintf("%s is a vpn/proxy worker application that sets up everything for you and connects to a central control server", AppName)),
		kong.DefaultEnvars("UNCHAINED_WORKER"),
		kong.Vars{
			"appname":       AppName,
			"dns":           defaults.Dns,
			"log_level":     defaults.LogLevel,
			"protos":        strings.Join(protocols.List, ","),
			"default_proto": protocols.Vless,
		},
	)
	err := appctx.Run()
	appctx.FatalIfErrorf(err)
}
