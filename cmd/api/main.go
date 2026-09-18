package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/getshelf/backend/internal/adapters/httpapi"
	"github.com/getshelf/backend/internal/adapters/postgres"
	"github.com/getshelf/backend/internal/config"
	"github.com/getshelf/backend/internal/modules/account"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()

	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	)

	appConfig, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New(
		"file://db/migrations",
		appConfig.DB.URL(),
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	db, err := postgres.Open(
		ctx,
		appConfig.DB,
	)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	accountStore := postgres.NewAccountStore(db)
	accounts := account.NewService(accountStore)

	handler := httpapi.New(httpapi.Dependencies{
		Accounts: accounts,
		Logger: logger,
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	logger.Info(
		"HTTP server starting",
		"address",
		server.Addr,
	)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
