package rest

import (
	"context"
	"net/http"

	"go.uber.org/fx"

	"github.com/getshelf/backend/internal/config"
)

func NewHTTPServer(lifecycle fx.Lifecycle, cfg config.Config, router http.Handler) *http.Server {
	server := &http.Server{Addr: cfg.HTTP.Address(), Handler: router}
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() { _ = server.ListenAndServe() }()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
	return server
}
