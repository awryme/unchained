package main

import "github.com/awryme/unchained/appconfig"

type RuntimeParams struct {
	LogLevel string   `help:"sing-box log level" default:"${log_level}"`
	DNS      string   `help:"dns address, in sing-box format" default:"${dns}"`
	Proto    string   `short:"p" help:"set used protocol: ${enum}" enum:"${protos}" default:"${default_proto}"`
	ID       string   `help:"proxy id (used to identify proxy in client apps), random by default"`
	Tags     []string `help:"proxy tags (used to identify proxy in client apps)"`
}

func (rp *RuntimeParams) GetRuntimeParams() *appconfig.RuntimeParams {
	return &appconfig.RuntimeParams{
		LogLevel: rp.LogLevel,
		DNS:      rp.DNS,
		ID:       rp.ID,
		Proto:    rp.Proto,
		Tags:     rp.Tags,
	}
}
