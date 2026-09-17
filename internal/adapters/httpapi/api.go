package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/getshelf/backend/internal/modules/account"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	Accounts *account.Service
	Logger   *slog.Logger
}

func New(dependencies Dependencies) http.Handler {
	if dependencies.Accounts == nil {
		panic("Accounts is required")
	}
	if dependencies.Logger == nil {
		panic("Logger is required")
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	config := huma.DefaultConfig(
		"Shelf API",
		"1.0.0",
	)

	api := humachi.New(router, config)

	v1 := huma.NewGroup(api, "/api/v1")

	registerAccountRoutes(
		v1,
		dependencies.Accounts,
		dependencies.Logger,
	)

	return router
}
