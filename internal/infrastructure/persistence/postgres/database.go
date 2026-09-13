package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/getshelf/backend/internal/config"
)

func NewDatabase(cfg config.Config) (*sqlx.DB, error) {
	return sqlx.Open("postgres", cfg.DB.URL())
}
