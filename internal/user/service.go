package user

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielbukowski/recipe-app-backend/gen/sqlc"
	"github.com/danielbukowski/recipe-app-backend/internal/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const queryExecutionTimeout = 3 * time.Second
const acquireConnectionTimeout = 3 * time.Second

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

	connCtx, cancelConnCtx := context.WithTimeout(ctx, acquireConnectionTimeout)
	defer cancelConnCtx()

	err = s.dbpool.AcquireFunc(connCtx, func(c *pgxpool.Conn) error {
		qCtx, cancelQCtx := context.WithTimeout(ctx, queryExecutionTimeout)
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
	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr):
			switch pgErr.Code {
			case "23505":
				return echo.NewHTTPError(http.StatusBadRequest, "user with this email already exists")
			default:
				s.logger.Error("createUser method got uncaught database error", zap.Error(err))
				return err
			}
		case errors.Is(err, context.DeadlineExceeded):
			return echo.NewHTTPError(http.StatusRequestTimeout)
		default:
			s.logger.Error("createUser method got uncaught error", zap.Error(err))
			return err
		}
	}

	return err
}
