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
	clilog.Log("Log level:", cfg.LogLevel)
	clilog.Log("DNS:", cfg.DNS)
	clilog.Log("Protocol:", cfg.Proto)
	clilog.Log("Port:", cfg.Listen.Port)
	clilog.Log("Reality server:", cfg.Reality.Server)
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
