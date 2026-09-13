package routers

import (
	"encoding/json"
	"errors"
	"net/http"

	applicationUser "github.com/getshelf/backend/internal/application/user"
	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
	"github.com/getshelf/backend/internal/presentation/rest/schemas"
)

type RegisterHandler struct {
	service applicationUser.RegistrationService
}

func NewRegisterHandler(service applicationUser.RegistrationService) RegisterHandler {
	return RegisterHandler{service: service}
}

func (handler RegisterHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeError(responseWriter, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input schemas.RegisterRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result, err := handler.service.Register(request.Context(), applicationUser.RegisterRequest{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		writeDomainError(responseWriter, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(responseWriter).Encode(schemas.RegisterResponse{ID: result.ID.String()})
}

func writeDomainError(responseWriter http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, domainErrors.ErrInvalidUsername), errors.Is(err, domainErrors.ErrInvalidEmail), errors.Is(err, domainErrors.ErrPasswordTooWeak):
		status = http.StatusBadRequest
	case errors.Is(err, domainErrors.ErrUsernameTaken), errors.Is(err, domainErrors.ErrEmailTaken):
		status = http.StatusConflict
	}
	writeError(responseWriter, status, err.Error())
}

func writeError(responseWriter http.ResponseWriter, status int, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)
	_ = json.NewEncoder(responseWriter).Encode(schemas.ErrorResponse{Error: message})
}
