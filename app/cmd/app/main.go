package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"vox/internal/config"
	"vox/internal/redis"
	"vox/internal/server"
)

func main() {
	cfg := config.New()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.GetLogLevel()}))
	slog.SetDefault(logger)

	rdb := redis.Init(cfg)
	srv := server.New(cfg, rdb)
	mux := srv.CreateHandlers()

	slog.Info("cmd/app/main.go", "Message", "Starting server", "port", cfg.GetPort(), "instance", cfg.GetInstanceName())
	log.Fatal(http.ListenAndServe(":"+cfg.GetPort(), mux))
}
