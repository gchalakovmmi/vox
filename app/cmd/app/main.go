package main

import (
	"os"
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

	slog.Info("cmd/app/main.go", "Message", "Extracted Environment Variables", "PORT", cfg.GetPort(), "INSTANCE_NAME", cfg.GetInstanceName(), "LOG_LEVEL", cfg.GetLogLevel())

	srv := server.New(cfg/*, db*/, server.Routes)
	mux := srv.CreateHandlers()

	slog.Info("cmd/app/main.go", "Message", "Starting HTTP server", "Port", cfg.GetPort())
	log.Fatal(http.ListenAndServe(":"+cfg.GetPort(), mux))
}
