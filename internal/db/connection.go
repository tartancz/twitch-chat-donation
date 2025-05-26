package db

import (
	"TwitchDonoCalculator/internal/config"
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(cfg config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
