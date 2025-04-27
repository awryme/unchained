package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/awryme/unchained/unchained/defaults"
	"github.com/awryme/unchained/unchained/protocols"
)

const AppName = "unchained"
const ConfigName = AppName + ".json"

type App struct {
	Run   CmdRun   `cmd:"" help:"run vpn server, generate config if it doesn't exist"`
	Print CmdPrint `cmd:"" help:"print connection info for client"`
}

func main() {
	printVersion()
	var app App
	appctx := kong.Parse(&app,
		kong.Name(AppName),
		kong.Description(fmt.Sprintf("%s is a vpn/proxy application that sets up everything for you", AppName)),
		kong.DefaultEnvars("UNCHAINED"),
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
