package middleware

import (
	"database/sql"
	"os"
	"fmt"
	// "net/http"
	// "log/slog"
	_ "github.com/lib/pq"
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
