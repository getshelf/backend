package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	v1 "github.com/getshelf/backend/internal/handler/rest/v1"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func RootRouter() http.Handler {
	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("Shelf API", "1.0.0"))

	public := huma.NewGroup(api, "")
	protected := huma.NewGroup(api, "")

	v1.RegisterRoutes(public, protected)

	return router
}
