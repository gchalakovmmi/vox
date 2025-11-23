package postgres

import (
	"database/sql"
	"os"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"vox/internal/config"
)

const (
    PGErrUniqueViolation = "23505"
)

func WithDBConn(cfg *config.Config, next func(*sql.DB) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := sql.Open("postgres",
			fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
				cfg.GetPostgresUser(),
				cfg.GetPostgresPassword(),
				cfg.GetPostgresDB(),
				cfg.GetPostgresHost()))
		if err != nil {
			slog.Error("internal/middleware/postgres.go", "Message", "Connecting to DB failed", "Error", err)
			http.Error(w, "Database connection failed", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		if err := conn.Ping(); err != nil {
			slog.Error("internal/middleware/postgres.go", "Message", "Connecting to DB failed", "Error", err)
			http.Error(w, "Database ping failed", http.StatusInternalServerError)
			return
		}

		next(conn).ServeHTTP(w, r)
	})
}
