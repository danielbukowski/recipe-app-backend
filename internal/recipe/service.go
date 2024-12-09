package recipe

import (
	"context"
	"errors"
	"time"

	"github.com/danielbukowski/recipe-app-backend/gen/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const databaseConnectionTimeout = 4 * time.Second

type service struct {
	Logger *zap.Logger
	dbpool *pgxpool.Pool
}

func NewService(logger *zap.Logger, dbpool *pgxpool.Pool) *service {
	return &service{
		Logger: logger,
		dbpool: dbpool,
	}
}

func (s *service) GetRecipeById(ctx context.Context, recipeId uuid.UUID) (sqlc.Recipe, error) {
	poolCtx, cancel := context.WithTimeout(ctx, databaseConnectionTimeout)
	defer cancel()

	conn, err := s.dbpool.Acquire(poolCtx)
	if err != nil {
		return sqlc.Recipe{}, errors.New("failed to connect to database")
	}
	defer conn.Release()

	q := sqlc.New(conn)

	dbCtx, cancel := context.WithTimeout(ctx, databaseConnectionTimeout)
	defer cancel()

	return q.GetRecipeById(dbCtx, recipeId)
}

func (s *service) deleteRecipeById(ctx context.Context, recipeid uuid.UUID) error {
	poolCtx, cancel := context.WithTimeout(ctx, databaseConnectionTimeout)
	defer cancel()

	conn, err := s.dbpool.Acquire(poolCtx)
	if err != nil {
		return errors.New("failed to connect to database")
	}
	defer conn.Release()

	q := sqlc.New(conn)

	dbCtx, cancel := context.WithTimeout(ctx, databaseConnectionTimeout)
	defer cancel()

	return q.DeleteRecipeById(dbCtx, recipeid)
}

