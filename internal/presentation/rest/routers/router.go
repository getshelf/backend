package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(registerHandler RegisterHandler, authHandler AuthHandler, authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	router.Mount("/api", NewAPIRouter(registerHandler, authHandler, authMiddleware))
	return router
}

func NewAPIRouter(registerHandler RegisterHandler, authHandler AuthHandler, authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	router.Mount("/v1", NewV1Router(registerHandler, authHandler, authMiddleware))
	return router
}

func NewV1Router(registerHandler RegisterHandler, authHandler AuthHandler, authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	router.Mount("/users", NewUsersRouter(registerHandler, authHandler, authMiddleware))
	return router
}

func NewUsersRouter(registerHandler RegisterHandler, authHandler AuthHandler, authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	registerUserRoutes(router, registerHandler, authHandler, authMiddleware)
	return router
}

func registerUserRoutes(router chi.Router, registerHandler RegisterHandler, authHandler AuthHandler, authMiddleware func(http.Handler) http.Handler) {
	router.Post("/register", registerHandler.ServeHTTP)
	router.Post("/login", authHandler.Login)
	router.Post("/refresh", authHandler.Refresh)
	router.Group(func(protected chi.Router) {
		protected.Use(authMiddleware)
		protected.Post("/logout", authHandler.Logout)
		protected.Get("/profile", authHandler.Profile)
		protected.Patch("/profile", authHandler.PatchProfile)
	})
}
