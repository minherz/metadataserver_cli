package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/minherz/metadataserver"
)

const undefinedStringFlag = "_undefined_"

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

	file := flag.String("config", undefinedStringFlag, "Path to configuration file")
	address := flag.String("address", metadataserver.DefaultAddress, "Serving address")
	port := flag.Int("port", metadataserver.DefaultPort, "Port to listen for requests")

	flag.Parse()

	ops := []metadataserver.Option{}
	if *file != undefinedStringFlag {
		ops = append(ops, metadataserver.WithConfigFile(*file))
	} else {
		ops = append(ops, metadataserver.WithAddress(*address), metadataserver.WithPort(*port))
	}
	srv, err := metadataserver.New(ops...)
	if err != nil {
		slog.Error("failed to create new server", slog.String("error", err.Error()))
		return err
	}
	srv.Start()

	// wait for interrupt and gracefully stop web server after 5 sec
	<-ctx.Done()
	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 5*time.Second)
	defer cancel()
	return srv.Stop(shutdownCtx)
}
