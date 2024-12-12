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

func (s *service) getRecipeById(ctx context.Context, recipeId uuid.UUID) (dto recipeResponse, err error) {
	var dbRecipe sqlc.Recipe

	err = s.dbpool.AcquireFunc(ctx, func(c *pgxpool.Conn) error {
	dbCtx, cancel := context.WithTimeout(ctx, queryExecutionTimeout)
	defer cancel()

	q := sqlc.New(c)

		dbRecipe, err = q.GetRecipeById(dbCtx, recipeId)
		return err
	})
	if err != nil {
		return dto, err
	}

	dto = recipeResponse{
		Title:     dbRecipe.Title,
		Content:   dbRecipe.Content,
		CreatedAt: dbRecipe.CreatedAt.Time,
		UpdatedAt: dbRecipe.CreatedAt.Time,
	}

	return dto, err
}

func (s *service) deleteRecipeById(ctx context.Context, recipeID uuid.UUID) error {
	connCtx, cancel := context.WithTimeout(ctx, acquireConnectionTimeout)
	defer cancel()

	return s.dbpool.AcquireFunc(connCtx, func(c *pgxpool.Conn) error {
		q := sqlc.New(c)

		qCtx, cancel := context.WithTimeout(ctx, queryExecutionTimeout)
		defer cancel()

		return q.DeleteRecipeById(qCtx, recipeID)
	})

}

func (s *service) CreateNewRecipe(ctx context.Context, recipe newRecipeRequest) error {
	connCtx, cancel := context.WithTimeout(ctx, acquireConnectionTimeout)
	defer cancel()

	id, err := uuid.NewV7()
	if err != nil {
		return errors.Join(errors.New("failed to generate UUID"), err)
	}

	return s.dbpool.AcquireFunc(connCtx, func(c *pgxpool.Conn) error {
		qCtx, cancel := context.WithTimeout(ctx, queryExecutionTimeout)
		defer cancel()

		q := sqlc.New(c)

		return q.CreateRecipe(
			qCtx,
			sqlc.CreateRecipeParams{
				RecipeID: id,
				Title:    recipe.title,
				Content:  recipe.content,
			},
		)
	})
}
