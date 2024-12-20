package recipe

import (
	"context"
	"errors"
	"time"

	"github.com/danielbukowski/recipe-app-backend/gen/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const queryExecutionTimeout = 4 * time.Second
const acquireConnectionTimeout = 3 * time.Second

type service struct {
	logger *zap.Logger
	dbpool *pgxpool.Pool
}

func NewService(logger *zap.Logger, dbpool *pgxpool.Pool) *service {
	return &service{
		logger: logger,
		dbpool: dbpool,
	}
}

func (s *service) GetRecipeById(ctx context.Context, recipeId uuid.UUID) (r sqlc.Recipe, err error) {
	err = s.dbpool.AcquireFunc(ctx, func(c *pgxpool.Conn) error {
		dbCtx, cancelDbCtx := context.WithTimeout(ctx, queryExecutionTimeout)
		defer cancelDbCtx()

		q := sqlc.New(c)

		r, err = q.GetRecipeById(dbCtx, recipeId)
		return err
	})
	if err != nil {
		return r, err
	}

	return r, err
}

func (s *service) DeleteRecipeById(ctx context.Context, recipeID uuid.UUID) error {
	connCtx, cancelConnCtx := context.WithTimeout(ctx, acquireConnectionTimeout)
	defer cancelConnCtx()

	return s.dbpool.AcquireFunc(connCtx, func(c *pgxpool.Conn) error {
		q := sqlc.New(c)

		qCtx, cancelQCtx := context.WithTimeout(ctx, queryExecutionTimeout)
		defer cancelQCtx()

		return q.DeleteRecipeById(qCtx, recipeID)
	})

}

func (s *service) CreateNewRecipe(ctx context.Context, newRecipeRequest NewRecipeRequest) (uuid.UUID, error) {
	var id uuid.UUID

	connCtx, cancelConnCtx := context.WithTimeout(ctx, acquireConnectionTimeout)
	defer cancelConnCtx()

	id, err := uuid.NewV7()
	if err != nil {
		return id, errors.Join(errors.New("failed to generate UUID"), err)
	}

	err = s.dbpool.AcquireFunc(connCtx, func(c *pgxpool.Conn) error {
		qCtx, cancelQCtx := context.WithTimeout(ctx, queryExecutionTimeout)
		defer cancelQCtx()

		q := sqlc.New(c)

		id, err = q.CreateRecipe(
			qCtx,
			sqlc.CreateRecipeParams{
				RecipeID: id,
				Title:    newRecipeRequest.Title,
				Content:  newRecipeRequest.Content,
			},
		)
		return err
	})

	return id, err
}

func (s *service) UpdateRecipeById(ctx context.Context, id uuid.UUID, updatedAt pgtype.Timestamp, newRecipeRequest UpdateRecipeRequest) error {
	connCtx, cancelConnCtx := context.WithTimeout(ctx, acquireConnectionTimeout)
	defer cancelConnCtx()

	tx, err := s.dbpool.Begin(connCtx)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)

	qCtx, cancelQCtx := context.WithTimeout(ctx, queryExecutionTimeout)
	defer cancelQCtx()

	err = q.UpdateRecipeById(qCtx, sqlc.UpdateRecipeByIdParams{
		RecipeID:  id,
		UpdatedAt: updatedAt,
		Title:     newRecipeRequest.Title,
		Content:   newRecipeRequest.Content,
		NewUpdatedAt: pgtype.Timestamp{
			Time:             time.Now(),
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
	})
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	return errors.Join(err, tx.Commit(ctx))
}
