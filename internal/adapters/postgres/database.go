package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/getshelf/backend/internal/config"
	_ "github.com/lib/pq"
)

func Open(
	ctx context.Context,
	config config.DBConfig,
) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		config.URL(),
	)

	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	pingContext, cancel := context.WithTimeout(
		ctx,
		5 * time.Second,
	)

	defer cancel()

	if err := db.PingContext(pingContext); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}
