package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getshelf/backend/internal/adapters/httpapi"
	"github.com/getshelf/backend/internal/adapters/postgres"
	"github.com/getshelf/backend/internal/adapters/webpage"
	"github.com/getshelf/backend/internal/config"
	"github.com/getshelf/backend/internal/modules/account"
	"github.com/getshelf/backend/internal/modules/collection"
	"github.com/getshelf/backend/internal/modules/item"
	"github.com/getshelf/backend/internal/modules/session"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

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

	m.Log = migrationLogger{logger: logger}

	if err := m.Up(); err != nil  {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("database migrations already up to date")
		} else {
			log.Fatal(err)
		}
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

	sessionStore := postgres.NewSessionStore(db)
	sessions := session.NewService(sessionStore, time.Hour * 24 * 14)

	collectionStore := postgres.NewCollectionStore(db)
	collections := collection.NewService(collectionStore)

	itemStore := postgres.NewItemStore(db)
	items := item.NewService(itemStore)

	extractor := webpage.NewExtractor()

	extractionWorker := item.NewExtractionWorker(
		itemStore,
		extractor,
		logger,
	)


	handler := httpapi.New(httpapi.Dependencies{
		Accounts: accounts,
		Logger: logger,
		Sessions: sessions,
		Collections: collections,
		Items: items,
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	workerDone := make(chan struct{})

	go func() {
		defer close(workerDone)
		extractionWorker.Run(ctx)
	}()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info(
			"HTTP server starting",
			"address", server.Addr,
		)
		serverErrors <- server.ListenAndServe()
	}()

	select {
		case <- ctx.Done():
			// Shutdown signal received
		case err := <- serverErrors:
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("HTTP Server failed", "error", err)
			}
	}

	stop()

	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		10 * time.Second,
	)

	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown failed", "error", err)
		_  = server.Close()
	}

	// Wait before the existing deferred db.Close() runs.
	<-workerDone
}
