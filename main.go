package main

import (
	"context"
	_ "embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gabe565.com/trmnl-nightscout/internal/config"
	"gabe565.com/trmnl-nightscout/internal/server"
	"gabe565.com/trmnl-nightscout/internal/util"
)

//go:generate go run ./internal/generate/docs

var version = "beta"

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	conf.Version = version

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	slog.Info("TRMNL Nightscout", "version", version, "commit", util.GetCommit())

	return server.New(conf).ListenAndServe(ctx)
}
