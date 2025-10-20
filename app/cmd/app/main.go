package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"log/slog"
	"vox/internal/config"
	"vox/internal/server"
)

func main() {
	cfg := config.New()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	slog.Info("Environment Variables", "PORT", cfg.Port, "INSTANCE_NAME", cfg.InstanceName, "LOG_LEVEL", cfg.LogLevel)

	srv := server.New(cfg/*, db*/)
	mux := srv.CreateHandlers()

	slog.Info(fmt.Sprintf("Starting HTTP server on port %s", cfg.Port))
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
