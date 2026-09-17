package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/getshelf/backend/internal/config"
	"github.com/getshelf/backend/internal/adapters/httpapi"
	"github.com/getshelf/backend/internal/adapters/postgres"
	"github.com/getshelf/backend/internal/modules/account"
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
