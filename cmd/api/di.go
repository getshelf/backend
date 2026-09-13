package main

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	applicationUser "github.com/getshelf/backend/internal/application/user"
	"github.com/getshelf/backend/internal/config"
	"github.com/getshelf/backend/internal/domain/user/ports"
	userService "github.com/getshelf/backend/internal/domain/user/services"
	postgres "github.com/getshelf/backend/internal/infrastructure/persistence/postgres"
	"github.com/getshelf/backend/internal/infrastructure/persistence/postgres/repositories"
	infraServices "github.com/getshelf/backend/internal/infrastructure/services"
)

func di() fx.Option {
	return fx.Options(
		fx.Provide(
			config.Load,
			postgres.NewDatabase,
			func(db *sqlx.DB) *sql.DB { return db.DB },
			repositories.NewUserRepository,
			postgres.NewTransactionManager,
			fx.Annotate(
				func() infraServices.BcryptPasswordHasher { return infraServices.NewBcryptPasswordHasher(0) },
				fx.As(new(ports.PasswordHasher)),
			),
			func() ports.PasswordVerifier { return infraServices.NewBcryptPasswordHasher(0) },
			fx.Annotate(infraServices.NewJWTService, fx.As(new(ports.TokenService))),
			fx.Annotate(infraServices.NewRedisBlacklist, fx.As(new(ports.TokenBlacklist))),
			func() ports.IDGenerator { return infraServices.RandomUserIDGenerator{} },
			func() ports.Clock { return infraServices.SystemClock{} },
			func(repository repositories.UserRepository) ports.UniquenessChecker { return repository },
			func(repository repositories.UserRepository) ports.UserRepository { return repository },
			func(repository repositories.UserRepository) ports.AuthRepository { return repository },
			func(manager postgres.TransactionManager) ports.TxManager { return manager },
			func() ports.EventPublisher { return nil },
			userService.NewRegistrationService,
			applicationUser.NewRegistrationService,
			applicationUser.NewAuthService,
		),
		fx.Invoke(func(lifecycle fx.Lifecycle, db *sql.DB) {
			lifecycle.Append(fx.Hook{OnStart: func(context.Context) error { return postgres.RunMigrations(db) }})
		}),
	)
}
