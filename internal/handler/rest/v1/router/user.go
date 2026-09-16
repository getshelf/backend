package router

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/getshelf/backend/internal/handler/rest/v1/schema"
)

func RegisterUserGroup(public huma.API, protected huma.API) {
	publicGroup := huma.NewGroup(public, "/users")
	protectedGroup := huma.NewGroup(protected, "/users")

	RegisterUserRoutes(publicGroup, protectedGroup)
}

func RegisterUserRoutes(public huma.API, protected huma.API) {
	huma.Register(public, huma.Operation{
		Method: http.MethodPost,
		Path:	"",
		Summary: "Create User",
		Description: "Create a new user",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, request *schema.CreateUserRequest) (*schema.CreateUserResponse, error) {
		return &schema.CreateUserResponse{
			Body: schema.CreateUserResponseBody{
				ID: "12",
			},
		}, nil
	})
}
