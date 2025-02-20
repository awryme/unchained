package main

import (
	"fmt"
	"os"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/trojan"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/file"
)

func printInfo(cfg appconfig.Config) error {
	fmt.Println("Sing-box log_level:", cfg.LogLevel)
	fmt.Println("Sing-box DNS:", cfg.DNS)
	fmt.Println("Reality server:", cfg.RealityServer)
	fmt.Println("Reality max time diff:", cfg.RealityTimeDiff)

	url := trojan.MakeUrl(cfg)
	fmt.Println("URL:", url)

	qr, err := qrcode.New(url)
	if err != nil {
		return fmt.Errorf("creating qr code: %w", err)
	}
	qr.Dimension()
	tw := file.New(os.Stdout)
	err = qr.Save(tw)
	if err != nil {
		return fmt.Errorf("printing qr code: %w", err)
	}
	return nil
}
