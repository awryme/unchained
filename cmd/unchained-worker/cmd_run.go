package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/awryme/unchained/app/appconfig"
	"github.com/awryme/unchained/app/clilog"
	"github.com/awryme/unchained/app/singbox/memoryuserstore"
	"github.com/awryme/unchained/app/singbox/singboxserver"
)

type CmdRun struct {
	LogLevel string   `help:"sing-box log level" default:"${log_level}"`
	DNS      string   `help:"dns address, in sing-box format" default:"${dns}"`
	ID       string   `help:"proxy id (used to identify proxy in client apps), random by default"`
	Tags     []string `help:"proxy tags (used to identify proxy in client apps)"`

	ControlAddr *url.URL `help:"control api address" required:""`
}

func (c *CmdRun) Run(app *App) error {
	ctx := context.Background()

	cfg, err := c.getConfig(ctx, app.Config)
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	inbounds, err := c.getInbounds(cfg)
	if err != nil {
		return fmt.Errorf("create inbounds: %w", err)
	}

	instance, err := singboxserver.Run(ctx, cfg.Singbox, inbounds...)
	if err != nil {
		return fmt.Errorf("run singbox server: %w", err)
	}
	clilog.Log("Started singbox at", time.Now().Format(time.DateTime))

	// register worker with a wrapper from unchained-control
	// api, err := workerapi.New(workerapi.Params{
	// 	PublicIP:       cfg.AppInfo.PublicIP,
	// 	WorkerID:       cfg.Worker.ID,
	// 	Listen:         cfg.Worker.Listen,
	// 	ControlAddr:    c.ControlAddr,
	// 	JwtSecret:      cfg.Worker.JwtSecret,
	// 	EncodedCert:    cfg.Worker.EncodedCert,
	// 	EncodedCertKey: cfg.Worker.EncodedCertKey,
	// })
	// if err != nil {
	// 	return fmt.Errorf("create worker api: %w", err)
	// }

	// if err := api.Register(); err != nil {
	// 	return fmt.Errorf("register worker: %w", err)
	// }

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	clilog.Log("Got ctrl+c / interrupt, quitting")

	return errors.Join(
		// shutdown api
		// api.Shutdown(ctx),
		instance.Close(),
	)
}

func (rp *CmdRun) GetRuntimeParams() *appconfig.UnchainedWorkerRuntimeParams {
	return &appconfig.UnchainedWorkerRuntimeParams{
		LogLevel: rp.LogLevel,
		DNS:      rp.DNS,
		Tags:     rp.Tags,
	}
}

func (c *CmdRun) getConfig(ctx context.Context, file string) (cfg appconfig.UnchainedWorker, err error) {
	params := c.GetRuntimeParams()

	cfg, err = appconfig.ReadUnchainedWorker(file, params)
	if errors.Is(err, os.ErrNotExist) {
		// cfg file not found, generate new one
		err = cfg.Generate(ctx, params)
	}
	if err != nil {
		return cfg, err
	}

	// write any changes that we applied
	err = appconfig.WriteUnchainedWorker(cfg, file)
	return cfg, err
}

func (c *CmdRun) getInbounds(cfg appconfig.UnchainedWorker) ([]singboxserver.InboundMaker, error) {
	trojanUserStore := memoryuserstore.NewTrojan()
	trojanInbound := singboxserver.NewInboundTrojan(cfg.ListenTrojan, cfg.Singbox, trojanUserStore)

	vlessUserStore := memoryuserstore.NewVless()
	vlessInbound := singboxserver.NewInboundVless(cfg.ListenVless, cfg.Singbox, vlessUserStore)

	return []singboxserver.InboundMaker{
		trojanInbound,
		vlessInbound,
	}, nil
}
