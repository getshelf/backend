package main

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/getshelf/backend/internal/presentation/rest"
	"github.com/getshelf/backend/internal/presentation/rest/middleware"
	"github.com/getshelf/backend/internal/presentation/rest/routers"
)

func server() fx.Option {
	return fx.Options(
		fx.Provide(
			routers.NewRegisterHandler,
			routers.NewAuthHandler,
			middleware.NewAuthMiddleware,
			fx.Annotate(routers.NewRouter, fx.As(new(http.Handler))),
			rest.NewHTTPServer,
		),
		fx.Invoke(func(*http.Server) {}),
	)
}
