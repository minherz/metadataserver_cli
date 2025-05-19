package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/minherz/metadataserver"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.Error("invalid run termination", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	ops, err := ConfigOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup server: %s", err.Error())
		flag.Usage()
		return nil
	}
	srv, err := metadataserver.New(ops...)
	if err != nil {
		slog.Error("failed to create new server", slog.String("error", err.Error()))
		return err
	}
	if err := srv.Start(context.Background()); err != nil {
		return err
	}
	// wait for interrupt and gracefully stop web server after 5 sec
	<-ctx.Done()
	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 5*time.Second)
	defer cancel()
	return srv.Stop(shutdownCtx)
}
