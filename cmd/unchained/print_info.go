package main

import (
	"fmt"

	"github.com/awryme/unchained/unchained/clilog"
	"github.com/awryme/unchained/unchained/config"
	"github.com/awryme/unchained/unchained/protocols"
	"github.com/awryme/unchained/unchained/urlmaker"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/file"
)

func printInfo(cfg config.Unchained, appInfo config.AppInfo, singbox config.Singbox) error {
	clilog.Log("log:", singbox.LogLevel)
	clilog.Log("dns:", singbox.DNS)
	clilog.Log("proto:", cfg.Proto)

	printProxyInfo(protocols.Vless, singbox.VlessProxy)
	printProxyInfo(protocols.Trojan, singbox.TrojanProxy)

	clilog.Log("name:", appInfo.Name(cfg.Proto))
	clilog.Log()

	urlMaker := urlmaker.New(appInfo, singbox)
	url, err := urlMaker.MakeByProto(cfg.Proto, cfg.VlessUUID, cfg.TrojanPassword)
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

func printProxyInfo(proto string, params config.ProxyParams) {
	clilog.Logf("%s params:", proto)
	clilog.Log("\tport:", params.Listen.Port())
	clilog.Log("\treality server:", params.Reality.Server)
}
