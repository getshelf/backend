package schema

import (
	"time"
)

type User struct {
	ID           string    `json:"id" readOnly:"true"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Body CreateUserRequestBody
}

type CreateUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Body struct {
		ID string `json:"id" doc:"Created user ID"`
	}
}

type CreateUserResponseBody struct {
	ID string `json:"id" doc:"Created user ID"`
}
