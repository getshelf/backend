package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	applicationUser "github.com/getshelf/backend/internal/application/user"
	"github.com/getshelf/backend/internal/domain/user/ports"
	userService "github.com/getshelf/backend/internal/domain/user/services"
	"github.com/getshelf/backend/internal/domain/user/values"
	postgres "github.com/getshelf/backend/internal/infrastructure/persistence/postgres"
	"github.com/getshelf/backend/internal/infrastructure/persistence/postgres/repositories"
	infraServices "github.com/getshelf/backend/internal/infrastructure/services"
	"github.com/getshelf/backend/internal/presentation/rest/routers"
)

type fixedIDGenerator struct{}

func (fixedIDGenerator) NewUserID(context.Context) (values.UserID, error) {
	return values.NewUserID("integration-user-1"), nil
}

type fixedClock struct{}

func (fixedClock) Now() time.Time {
	return time.Date(2026, 9, 13, 12, 0, 0, 0, time.UTC)
}

func TestRegisterHTTP(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := sqlx.NewDb(sqlDB, "sqlmock")
	repository := repositories.NewUserRepository(db)
	transactionManager := postgres.NewTransactionManager(db)
	domainService := userService.NewRegistrationService(
		repository,
		repository,
		infraServices.NewBcryptPasswordHasher(4),
		fixedIDGenerator{},
		fixedClock{},
		transactionManager,
		nil,
	)
	applicationService := applicationUser.NewRegistrationService(domainService)
	handler := routers.NewRouter(routers.NewRegisterHandler(applicationService), routers.NewAuthHandler(applicationUser.AuthService{}), func(next http.Handler) http.Handler { return next })

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`)).
		WithArgs("john_doe").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`)).
		WithArgs("john@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, username, email, password_hash, created_at) VALUES ($1, $2, $3, $4, $5)`)).
		WithArgs("integration-user-1", "john_doe", "john@example.com", sqlmock.AnyArg(), fixedClock{}.Now()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users/register", strings.NewReader(`{"username":"john_doe","email":"JOHN@EXAMPLE.COM","password":"StrongPass123"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code, "body: %s", response.Body.String())
	assert.Contains(t, response.Body.String(), `"id":"integration-user-1"`)
	require.NoError(t, mock.ExpectationsWereMet())
}

var _ ports.IDGenerator = fixedIDGenerator{}
var _ ports.Clock = fixedClock{}
