package main

import (
	"fmt"

	"github.com/awryme/unchained/app/appconfig"
	"github.com/awryme/unchained/app/clilog"
	"github.com/awryme/unchained/app/singbox/singboxserver"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/file"
)

func printInfo(cfg appconfig.Unchained) error {
	clilog.Log("log:", cfg.Singbox.LogLevel)
	clilog.Log("dns:", cfg.Singbox.DNS)
	clilog.Log("proto:", cfg.Proto)
	clilog.Log("port:", cfg.Listen.Port())
	clilog.Log("reality_server:", cfg.Singbox.Reality.Server)
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
