package user

import (
	"context"
	"errors"
	"time"

	"github.com/danielbukowski/recipe-app-backend/gen/sqlc"
	"github.com/danielbukowski/recipe-app-backend/internal/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type service struct {
	logger         *zap.Logger
	dbpool         *pgxpool.Pool
	passwordHasher passwordHasher
}

type passwordHasher interface {
	CreateHashFromPassword(password string) (string, error)
}

func NewService(logger *zap.Logger, passwordHasher passwordHasher, dbppol *pgxpool.Pool) *service {
	return &service{
		logger:         logger,
		dbpool:         dbppol,
		passwordHasher: passwordHasher,
	}
}

func (s *service) CreateUser(ctx context.Context, user auth.SignUpRequest) error {
	hashedPassword, err := s.passwordHasher.CreateHashFromPassword(user.Password)
	if err != nil {
		return errors.Join(errors.New("failed to generate a hash from the password"), err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return errors.Join(errors.New("failed to generate UUID"), err)
	}

	connCtx, cancelConnCtx := context.WithTimeout(ctx, 3*time.Second)
	defer cancelConnCtx()

	err = s.dbpool.AcquireFunc(connCtx, func(c *pgxpool.Conn) error {
		qCtx, cancelQCtx := context.WithTimeout(ctx, 3*time.Second)
		defer cancelQCtx()

		q := sqlc.New(c)

		return q.CreateUser(qCtx,
			sqlc.CreateUserParams{
				UserID:   id,
				Email:    user.Email,
				Password: hashedPassword,
			},
		)
	})
	return err
}
