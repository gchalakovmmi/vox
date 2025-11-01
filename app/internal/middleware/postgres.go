package middleware

import (
	"database/sql"
	"os"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
)

const (
    PGErrUniqueViolation = "23505"
)

type DB struct {
	user			string
	password		string
	dbName			string
	host			string
	connectionString	string
	conn			*sql.DB
}

func NewDB() (*DB, error){
	db := DB{
		user:		os.Getenv("POSTGRES_USER"),
		password:	os.Getenv("POSTGRES_PASSWORD"),
		dbName:		os.Getenv("POSTGRES_DB"),
		host:		os.Getenv("POSTGRES_HOST"),
	}

	db.connectionString = fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		db.user, db.password, db.dbName, db.host)

	var err error
	db.conn, err = sql.Open("postgres", db.connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.conn.Ping(); err != nil {
		db.conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return &db, nil
}

func (db *DB) Conn() (*sql.DB){
	return db.conn
}

func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func WithDBConn(next func(*sql.DB) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, err := NewDB()
		if err != nil {
			slog.Error("internal/middleware/postgres.go", "Message", "Connecting to DB failed", "Error", err)
			http.Error(w, "Database connection failed", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		conn := db.Conn()
		slog.Debug("internal/middleware/postgres.go", "Message", "DB Connection Successful")
		
		next(conn).ServeHTTP(w, r)
	})
}
