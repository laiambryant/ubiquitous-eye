package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"ubiquitous-eye/packages/server"
	"ubiquitous-eye/packages/telemetry"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := telemetry.SetupOTelSDK(ctx)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	server.RunServer()

}
