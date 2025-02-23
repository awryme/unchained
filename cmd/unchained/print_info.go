package main

import (
	"fmt"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/clilog"
	"github.com/awryme/unchained/pkg/singboxserver"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/file"
)

func printInfo(cfg appconfig.Config) error {
	clilog.Log("log:", cfg.LogLevel)
	clilog.Log("dns:", cfg.DNS)
	clilog.Log("proto:", cfg.Proto)
	clilog.Log("port:", cfg.Listen.Port)
	clilog.Log("reality_server:", cfg.Reality.Server)
	clilog.Log("name:", cfg.Name())
	clilog.Log()

	url, err := singboxserver.MakeUrl(cfg)
	if err != nil {
		return fmt.Errorf("make proxy url: %w", err)
	}

	clilog.Log("Client URL:")
	clilog.Log(url)

	qr, err := qrcode.New(url)
	if err != nil {
		return fmt.Errorf("creating qr code: %w", err)
	}
	qr.Dimension()
	tw := file.New(clilog.Output)
	err = qr.Save(tw)
	if err != nil {
		return fmt.Errorf("printing qr code: %w", err)
	}
	return nil
}
