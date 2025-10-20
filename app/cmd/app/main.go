package main

import (
	"os"
	"log/slog"
	"vox/internal/config"
)

func main() {
	cfg := config.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)
	slog.Info("Environment Variables", "PORT", cfg.Port, "INSTANCE_NAME", cfg.InstanceName, "LOG_LEVEL", cfg.LogLevel)
}
