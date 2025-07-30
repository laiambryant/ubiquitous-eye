package main

import (
	"log/slog"
	"os"

	"github.com/laiambryant/ubiquitous-eye/packages/logger"
	"github.com/laiambryant/ubiquitous-eye/packages/services"
	"github.com/laiambryant/ubiquitous-eye/packages/utils"
)

func main() {
	logger.ConfigureLogger()
	err := services.CreateDeploySite(utils.GetDeployableSiteURI())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
