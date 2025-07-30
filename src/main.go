package main

import (
	"github.com/laiambryant/ubiquitous-eye/packages/logger"
	"github.com/laiambryant/ubiquitous-eye/packages/services"
	"github.com/laiambryant/ubiquitous-eye/packages/utils"
)

func main() {
	logger.ConfigureLogger()
	services.CreateDeploySite(utils.DEPLOYABLE_SITE_URI)
}
