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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.GetLogLevel()}))
	slog.SetDefault(logger)

	slog.Info("Environment Variables", "PORT", cfg.GetPort(), "INSTANCE_NAME", cfg.GetInstanceName(), "LOG_LEVEL", cfg.GetLogLevel())

	srv := server.New(cfg/*, db*/)
	mux := srv.CreateHandlers()

	slog.Info(fmt.Sprintf("Starting HTTP server on port %s", cfg.GetPort()))
	log.Fatal(http.ListenAndServe(":"+cfg.GetPort(), mux))
}
