package recipe

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielbukowski/recipe-app-backend/gen/sqlc"
	"github.com/danielbukowski/recipe-app-backend/internal/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const queryExecutionTimeout = 3 * time.Second
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

func (s *service) GetRecipeById(ctx context.Context, recipeId uuid.UUID) (RecipeResponse, error) {
	var recipeResponse RecipeResponse

	err := s.dbpool.AcquireFunc(ctx, func(c *pgxpool.Conn) error {
		dbCtx, cancelDbCtx := context.WithTimeout(ctx, queryExecutionTimeout)
		defer cancelDbCtx()

		q := sqlc.New(c)

		recipeFromDb, err := q.GetRecipeById(dbCtx, recipeId)
		if err != nil {
			return err
		}

		recipeResponse = RecipeResponse{
			Title:     recipeFromDb.Title,
			Content:   recipeFromDb.Content,
			CreatedAt: recipeFromDb.CreatedAt.Time,
			UpdatedAt: recipeFromDb.UpdatedAt.Time,
		}

		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return RecipeResponse{}, echo.NewHTTPError(http.StatusNotFound, shared.CommonResponse{Message: "could not find a recipe with this ID"})
		case errors.Is(err, context.DeadlineExceeded):
			return RecipeResponse{}, echo.NewHTTPError(http.StatusRequestTimeout)
		default:
			s.logger.Error("getRecipeById method got uncaught error", zap.Error(err))
			return RecipeResponse{}, err
		}
	}

	return recipeResponse, nil
}

func (s *service) DeleteRecipeById(ctx context.Context, recipeID uuid.UUID) error {
	connCtx, cancelConnCtx := context.WithTimeout(ctx, acquireConnectionTimeout)
	defer cancelConnCtx()

	err := s.dbpool.AcquireFunc(connCtx, func(c *pgxpool.Conn) error {
		q := sqlc.New(c)

		qCtx, cancelQCtx := context.WithTimeout(ctx, queryExecutionTimeout)
		defer cancelQCtx()

		return q.DeleteRecipeById(qCtx, recipeID)
	})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return echo.NewHTTPError(http.StatusNoContent)
		case errors.Is(err, context.DeadlineExceeded):
			return echo.NewHTTPError(http.StatusRequestTimeout)
		default:
			s.logger.Error("deleteRecipeById method got uncaught error", zap.Error(err))
			return err
		}
	}

	return nil
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
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return uuid.UUID{}, echo.NewHTTPError(http.StatusRequestTimeout)
		default:
			s.logger.Error("createNewRecipe method got uncaught error", zap.Error(err))
			return uuid.UUID{}, err
		}
	}

	return id, nil
}

func (s *service) UpdateRecipeById(ctx context.Context, id uuid.UUID, updatedAt time.Time, updateRecipeRequest UpdateRecipeRequest) error {
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
		RecipeID: id,
		UpdatedAt: pgtype.Timestamp{
			Time:             updatedAt,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
		Title:   updateRecipeRequest.Title,
		Content: updateRecipeRequest.Content,
		NewUpdatedAt: pgtype.Timestamp{
			Time:             time.Now(),
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
	})
	if err != nil {
		_ = tx.Rollback(ctx)

		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return echo.NewHTTPError(http.StatusConflict, shared.CommonResponse{Message: "conflict occurred when trying to update a recipe"})
		default:
			s.logger.Error("updateRecipeById method got uncaught error", zap.Error(err))
			return err
		}
	}

	return tx.Commit(ctx)
}
