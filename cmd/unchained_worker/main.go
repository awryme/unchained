package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/awryme/unchained/app/appconfig"
	"github.com/awryme/unchained/pkg/protocols"
)

const AppName = "unchained-worker"

type App struct {
	Config string `help:"file to store generated/edited config file" short:"c" default:"./${appname}.json"`
	Run    CmdRun `cmd:"" help:"run vpn worker, generate config if it doesn't exist"`
}

func main() {
	var app App
	appctx := kong.Parse(&app,
		kong.Name(AppName),
		kong.Description(fmt.Sprintf("%s is a vpn/proxy worker application that sets up everything for you and connects to a central control server", AppName)),
		kong.DefaultEnvars(""),
		kong.Vars{
			"appname":       AppName,
			"dns":           appconfig.DefaultDns,
			"log_level":     appconfig.DefaultLogLevel,
			"protos":        strings.Join(protocols.List, ","),
			"default_proto": protocols.Vless,
		},
	)
	err := appctx.Run()
	appctx.FatalIfErrorf(err)
}
