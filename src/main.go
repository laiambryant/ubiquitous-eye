package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"ubiquitous-eye/packages/server"
	"ubiquitous-eye/packages/services"
	"ubiquitous-eye/packages/telemetry"
)

func main() {
	if os.Args != nil {
		if len(os.Args) > 1 {
			if os.Args[1] == "server-mode" {
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
		} else {
			print("Creating deploy site\n")
			services.CreateDeploySite()
			print("Deploy site created")
		}
	}

}
