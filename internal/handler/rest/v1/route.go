package v1

import (
	"github.com/danielgtaylor/huma/v2"
	router "github.com/getshelf/backend/internal/handler/rest/v1/router"
)

func RegisterGroup(public huma.API, protected huma.API) {
	publicGroup := huma.NewGroup(public, "/v1")
	protectedGroup := huma.NewGroup(protected, "/v1")

	RegisterRoutes(publicGroup, protectedGroup)
}

func RegisterRoutes(public huma.API, protected huma.API) {
	router.RegisterUserGroup(public, protected)
}
