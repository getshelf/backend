package user

import (
	"context"

	userService "github.com/getshelf/backend/internal/domain/user/services"
	"github.com/getshelf/backend/internal/domain/user/values"
)

type RegisterRequest struct {
	Email    string
	Password string
}

type RegisterResponse struct {
	ID values.UserID
}

type RegistrationService struct {
	domain userService.RegistrationService
}

func NewRegistrationService(domain userService.RegistrationService) RegistrationService {
	return RegistrationService{domain: domain}
}

func (service RegistrationService) Register(ctx context.Context, request RegisterRequest) (RegisterResponse, error) {
	id, err := service.domain.Register(ctx, request.Email, request.Password)
	if err != nil {
		return RegisterResponse{}, err
	}
	return RegisterResponse{ID: id}, nil
}
