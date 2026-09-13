package routers

import (
	"encoding/json"
	"net/http"

	applicationUser "github.com/getshelf/backend/internal/application/user"
	"github.com/getshelf/backend/internal/presentation/rest/middleware"
	"github.com/getshelf/backend/internal/presentation/rest/schemas"
)

type AuthHandler struct {
	service applicationUser.AuthService
}

func NewAuthHandler(service applicationUser.AuthService) AuthHandler {
	return AuthHandler{service: service}
}

func (handler AuthHandler) Login(responseWriter http.ResponseWriter, request *http.Request) {
	var input schemas.LoginRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid JSON body")
		return
	}
	pair, err := handler.service.Login(request.Context(), applicationUser.LoginRequest{Email: input.Email, Password: input.Password})
	if err != nil {
		middleware.WriteAuthError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, schemas.TokenResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken})
}

func (handler AuthHandler) Refresh(responseWriter http.ResponseWriter, request *http.Request) {
	var input schemas.RefreshRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid JSON body")
		return
	}
	pair, err := handler.service.Refresh(request.Context(), input.RefreshToken)
	if err != nil {
		middleware.WriteAuthError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, schemas.TokenResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken})
}

func (handler AuthHandler) Logout(responseWriter http.ResponseWriter, request *http.Request) {
	access := middleware.BearerToken(request)
	var input schemas.RefreshRequest
	_ = json.NewDecoder(request.Body).Decode(&input)
	if err := handler.service.Logout(request.Context(), access, input.RefreshToken); err != nil {
		writeError(responseWriter, http.StatusInternalServerError, "could not revoke tokens")
		return
	}
	responseWriter.WriteHeader(http.StatusNoContent)
}

func (handler AuthHandler) Profile(responseWriter http.ResponseWriter, request *http.Request) {
	profile, err := handler.service.Profile(request.Context(), middleware.UserID(request.Context()))
	if err != nil {
		middleware.WriteAuthError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, profileResponse(profile))
}

func (handler AuthHandler) PatchProfile(responseWriter http.ResponseWriter, request *http.Request) {
	var input schemas.UpdateProfileRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid JSON body")
		return
	}
	profile, err := handler.service.UpdateProfile(request.Context(), middleware.UserID(request.Context()), applicationUser.UpdateProfileRequest{Email: input.Email})
	if err != nil {
		middleware.WriteAuthError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, profileResponse(profile))
}

func profileResponse(profile applicationUser.Profile) schemas.ProfileResponse {
	return schemas.ProfileResponse{ID: profile.ID.String(), Email: profile.Email.String(), CreatedAt: profile.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")}
}

func writeJSON(responseWriter http.ResponseWriter, status int, body any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)
	_ = json.NewEncoder(responseWriter).Encode(body)
}
