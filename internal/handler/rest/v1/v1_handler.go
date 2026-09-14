package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func v1Router() http.Handler {
	r := chi.NewRouter()

	return r
}