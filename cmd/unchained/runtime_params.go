package main

import "github.com/awryme/unchained/app/appconfig"

type RuntimeParams struct {
	LogLevel string   `help:"sing-box log level" default:"${log_level}"`
	DNS      string   `help:"dns address, in sing-box format" default:"${dns}"`
	Proto    string   `short:"p" help:"set used protocol: ${enum}" enum:"${protos}" default:"${default_proto}"`
	Tags     []string `help:"proxy tags (used to identify proxy in client apps)"`
}

func (rp *RuntimeParams) GetRuntimeParams() *appconfig.UnchainedRuntimeParams {
	return &appconfig.UnchainedRuntimeParams{
		LogLevel: rp.LogLevel,
		DNS:      rp.DNS,
		Proto:    rp.Proto,
		Tags:     rp.Tags,
	}
}
