package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

type Config struct {
	DSN string
}

func NewDB(config Config) (*DB, error) {
	if config.DSN == "" {
		slog.Error("Database DSN is required")
		return nil, fmt.Errorf("database DSN is required")
	}

	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		slog.Error("Failed to open database connection",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Database connection established successfully")

	return &DB{
		DB: db,
	}, nil
}

func (db *DB) Close() error {
	slog.Info("Closing database connection")
	return db.DB.Close()
}

func (db *DB) Health() error {
	return db.Ping()
}
