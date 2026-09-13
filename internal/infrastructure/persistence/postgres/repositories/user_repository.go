package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
	"github.com/getshelf/backend/internal/domain/user/ports"
	"github.com/getshelf/backend/internal/domain/user/values"
	postgres "github.com/getshelf/backend/internal/infrastructure/persistence/postgres"
)

type UserRepository struct {
	db *sqlx.DB
}

type sqlExecutor interface {
	QueryRowxContext(context.Context, string, ...any) *sqlx.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepository{db: db}
}

func (repository UserRepository) IsUsernameTaken(ctx context.Context, username values.Username) (bool, error) {
	var exists bool
	err := repository.query(ctx).QueryRowxContext(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`, username.String()).Scan(&exists)
	return exists, err
}

func (repository UserRepository) IsEmailTaken(ctx context.Context, email values.Email) (bool, error) {
	var exists bool
	err := repository.query(ctx).QueryRowxContext(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`, email.String()).Scan(&exists)
	return exists, err
}

func (repository UserRepository) Save(ctx context.Context, user ports.User) error {
	_, err := repository.query(ctx).ExecContext(ctx, `INSERT INTO users (id, username, email, password_hash, created_at) VALUES ($1, $2, $3, $4, $5)`, user.ID().String(), user.Username().String(), user.Email().String(), user.PasswordHash().String(), user.CreatedAt())
	if err != nil {
		var postgresError *pq.Error
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			switch postgresError.Constraint {
			case "users_username_unique":
				return domainErrors.ErrUsernameTaken
			case "users_email_unique":
				return domainErrors.ErrEmailTaken
			}
		}
		return err
	}
	return nil
}

func (repository UserRepository) FindByLogin(ctx context.Context, login string) (ports.AuthUser, error) {
	return repository.findUser(ctx, `SELECT id, username, email, password_hash, created_at FROM users WHERE username = $1 OR email = $1 LIMIT 1`, login)
}

func (repository UserRepository) FindByID(ctx context.Context, id values.UserID) (ports.AuthUser, error) {
	return repository.findUser(ctx, `SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1`, id.String())
}

func (repository UserRepository) UpdateProfile(ctx context.Context, id values.UserID, username values.Username, email values.Email) error {
	_, err := repository.query(ctx).ExecContext(ctx, `UPDATE users SET username = $1, email = $2 WHERE id = $3`, username.String(), email.String(), id.String())
	return err
}

func (repository UserRepository) findUser(ctx context.Context, query string, args ...any) (ports.AuthUser, error) {
	var user ports.AuthUser
	err := repository.query(ctx).QueryRowxContext(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

func (repository UserRepository) query(ctx context.Context) sqlExecutor {
	if tx, ok := postgres.TxFromContext(ctx); ok {
		return tx
	}
	return repository.db
}

var _ ports.UniquenessChecker = UserRepository{}
var _ ports.UserRepository = UserRepository{}
var _ ports.AuthRepository = UserRepository{}
