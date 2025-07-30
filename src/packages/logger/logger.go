package logger

import (
	"log/slog"

	"github.com/laiambryant/ubiquitous-eye/configuration"
)

func ConfigureLogger() {
	var handler slog.Handler
	config, err := configuration.LoadConfig()
	if err != nil {
		logger := slog.New(handler).With(slog.String("level", "info"))
		slog.SetDefault(logger)
		logger.Warn("Default configuration loaded, error while loading configuration", "error", err.Error())
	}
	if config.File != "" {
		logger := slog.New(handler).With(slog.String("level", config.Level))
		slog.SetDefault(logger)
		logger.Info("Loaded configuration", "path", config.File, "level", config.Level)
	}
}
