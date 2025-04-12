package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/protocols"
)

const AppName = "unchained"

type App struct {
	Config   string      `help:"file to store generated/edited config file" short:"c" default:"./${appname}.json"`
	Run      CmdRun      `cmd:"" help:"run vpn server, generates config if it doesn't exist"`
	Print    CmdPrint    `cmd:"" help:"print connection info for client"`
	Generate CmdGenerate `cmd:"" help:"generate config without running the server"`
	Reset    CmdReset    `cmd:"" help:"cleans up configs/files used by this command"`
}

func main() {
	var app App
	appctx := kong.Parse(&app,
		kong.Name(AppName),
		kong.Description(fmt.Sprintf("%s is a vpn/proxy application that sets up everything for you", AppName)),
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
