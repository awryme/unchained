package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/awryme/unchained-control/unchained-control/workerapi"
	"github.com/awryme/unchained/unchained-worker/config"
	"github.com/awryme/unchained/unchained/clilog"
	"github.com/awryme/unchained/unchained/singbox/singboxserver"
	"github.com/awryme/unchained/unchained/urlmaker"
	"github.com/awryme/unchained/unchained/userstore"
	"github.com/awryme/unchained/unchained/workerstore"
)

type CmdRun struct {
	LogLevel string   `help:"sing-box log level" default:"${log_level}"`
	DNS      string   `help:"dns address, in sing-box format" default:"${dns}"`
	Tags     []string `help:"proxy tags (used to identify proxy in client apps)"`
	// todo: use id?
	ID string `help:"proxy id (used to identify proxy in client apps), random by default"`

	Control *url.URL `help:"control api address" required:""`

	Dir string `help:"dir to store config file and data" default:"./data/"`
}

func (cmd *CmdRun) Run(app *App) error {
	ctx := context.Background()

	cfg, err := cmd.getConfig(ctx)
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	userstore, err := userstore.NewFileStore(cmd.Dir)
	if err != nil {
		return fmt.Errorf("create userstore: %w", err)
	}

	inbounds, err := cmd.getInbounds(cfg, userstore)
	if err != nil {
		return fmt.Errorf("create inbounds: %w", err)
	}

	instance, err := singboxserver.Run(ctx, cfg.Singbox, inbounds...)
	if err != nil {
		return fmt.Errorf("run singbox server: %w", err)
	}
	clilog.Log("Started singbox at", time.Now().Format(time.DateTime))

	urlMaker := urlmaker.New(cfg.AppInfo, cfg.Singbox)

	workerstore, err := workerstore.NewFileStore(cmd.Dir, userstore, urlMaker)
	if err != nil {
		return fmt.Errorf("create workerstore: %w", err)
	}

	api, err := workerapi.New(cmd.Control, cfg.Worker.Listen, cfg.AppInfo.PublicIP, workerstore)
	if err != nil {
		return fmt.Errorf("create worker api: %w", err)
	}

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	clilog.Log("Got ctrl+c / interrupt, quitting")

	return errors.Join(
		// shutdown api
		api.Shutdown(time.Second*10),
		instance.Close(),
	)
}

func (cmd *CmdRun) getConfigName() (string, error) {
	if err := os.MkdirAll(cmd.Dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("make data dir: %w", err)
	}
	file := filepath.Join(cmd.Dir, ConfigName)
	return file, nil
}

func (cmd *CmdRun) getConfig(ctx context.Context) (cfg config.UnchainedWorker, err error) {
	params := &config.DynamicParams{
		LogLevel: cmd.LogLevel,
		DNS:      cmd.DNS,
		Tags:     cmd.Tags,
	}

	file, err := cmd.getConfigName()
	if err != nil {
		return cfg, fmt.Errorf("get config: %w", err)
	}

	cfg, err = config.Read(file, params)
	if errors.Is(err, os.ErrNotExist) {
		// cfg file not found, generate new one
		err = cfg.Generate(ctx, params)
	}
	if err != nil {
		return cfg, err
	}

	// write any changes that we applied
	err = config.Write(cfg, file)
	return cfg, err
}

func (cmd *CmdRun) getInbounds(cfg config.UnchainedWorker, store *userstore.FileStore) ([]singboxserver.InboundMaker, error) {
	trojanInbound := singboxserver.NewInboundTrojan(cfg.Singbox.TrojanProxy, store)
	vlessInbound := singboxserver.NewInboundVless(cfg.Singbox.VlessProxy, store)

	return []singboxserver.InboundMaker{
		trojanInbound,
		vlessInbound,
	}, nil
}
