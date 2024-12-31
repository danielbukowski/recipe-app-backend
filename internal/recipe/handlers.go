package recipe

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielbukowski/recipe-app-backend/gen/sqlc"
	"github.com/danielbukowski/recipe-app-backend/internal/shared"
	"github.com/danielbukowski/recipe-app-backend/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type handler struct {
	logger        *zap.Logger
	recipeService recipeService
}

type recipeService interface {
	GetRecipeById(context.Context, uuid.UUID) (sqlc.Recipe, error)
	DeleteRecipeById(context.Context, uuid.UUID) error
	CreateNewRecipe(context.Context, NewRecipeRequest) (uuid.UUID, error)
	UpdateRecipeById(context.Context, uuid.UUID, pgtype.Timestamp, UpdateRecipeRequest) error
}

func NewHandler(logger *zap.Logger, recipeService recipeService) *handler {
	return &handler{
		logger:        logger,
		recipeService: recipeService,
	}
}

//	@Summary		Create a recipe
//	@Description	Insert a new recipe by providing a request body with a title and a content for the recipe.
//
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			NewRecipeRequest	body		recipe.NewRecipeRequest	true	"Request body with title and content"
//
//	@Success		201					{object}	shared.CommonResponse
//	@Failure		400					{object}	shared.CommonResponse
//	@Failure		408					{object}	shared.CommonResponse
//	@Failure		500					{object}	shared.CommonResponse
//
//	@Router			/api/v1/recipes [POST]
func (h *handler) createRecipe(ctx *gin.Context) {
	var requestBody = NewRecipeRequest{}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
			"message": "missing JSON request body",
		})
		return
	}

	v := validator.New()

	if validateNewRecipeRequestBody(v, requestBody); !v.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request body did not pass the validation",
			"fields":  v.Errors,
		})
		return
	}

	recipeId, err := h.recipeService.CreateNewRecipe(ctx.Copy(), requestBody)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to save a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error(err.Error(),
				zap.Stack("stackError"),
				zap.String("recipeId", recipeId.String()),
				zap.String("method", ctx.Request.Method),
				zap.String("path", ctx.FullPath()),
			)
		}
		return
	}

	h.logger.Info("saved a new recipe to database")

	// TODO: find out a better way to get the address
	ctx.Header("Location", fmt.Sprintf("http://localhost:8080/api/v1/recipes/%v", recipeId.String()))
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "successfully saved a recipe",
	})
}

//	@Summary		Update a recipe
//	@Description	Update a title or a content of a recipe by ID.
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id					path		int							true	"ID for a recipe resource"
//	@Param			UpdateRecipeRequest	body		recipe.UpdateRecipeRequest	true	"Request body for updating title and content fields of a recipe"
//	@Success		204
//	@Failure		400					{object}	shared.CommonResponse
//	@Failure		408					{object}	shared.CommonResponse
//	@Failure		500					{object}	shared.CommonResponse
//	@Router			/api/v1/recipes/{id} [PUT]
func (h *handler) updateRecipeById(ctx *gin.Context) {
	recipeIdParam, ok := ctx.Params.Get("id")

	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "missing ID param for recipe",
		})
		return
	}

	recipeId, err := uuid.Parse(recipeIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "the received ID is not a valid UUID",
		})
		return
	}

	var requestBody = UpdateRecipeRequest{}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
			"message": "missing JSON request body",
		})
		return
	}

	recipeFromDb, err := h.recipeService.GetRecipeById(ctx, recipeId)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "could not find a recipe with this id",
			})
		case errors.Is(err, context.DeadlineExceeded):
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to update a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error(err.Error(),
				zap.Stack("stackError"),
				zap.String("recipeId", recipeId.String()),
				zap.String("method", ctx.Request.Method),
				zap.String("path", ctx.FullPath()),
			)
		}
		return
	}

	if requestBody.Title == "" {
		requestBody.Title = recipeFromDb.Title
	}

	if requestBody.Content == "" {
		requestBody.Content = recipeFromDb.Content
	}

	v := validator.New()

	if validateUpdateRecipeRequestBody(v, requestBody); !v.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request body did not pass the validation",
			"fields":  v.Errors,
		})
		return
	}

	err = h.recipeService.UpdateRecipeById(ctx.Copy(), recipeId, recipeFromDb.UpdatedAt, requestBody)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			ctx.JSON(http.StatusConflict, gin.H{
				"message": "conflict occurred when trying to update a recipe",
			})
		case errors.Is(err, context.DeadlineExceeded):
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to save a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error(err.Error(),
				zap.Stack("stackError"),
				zap.String("recipeId", recipeId.String()),
				zap.String("method", ctx.Request.Method),
				zap.String("path", ctx.FullPath()),
			)
		}
		return
	}

	h.logger.Info("updated a recipe",
		zap.String("recipeId", recipeId.String()),
	)

	ctx.Status(http.StatusNoContent)
}

func (h *handler) deleteRecipeById(ctx *gin.Context) {
	recipeIdParam, ok := ctx.Params.Get("id")

	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "missing ID param for recipe",
		})
		return
	}

	recipeId, err := uuid.Parse(recipeIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "the received ID is not a valid UUID",
		})
		return
	}

	err = h.recipeService.DeleteRecipeById(ctx.Copy(), recipeId)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			ctx.Status(http.StatusNoContent)
		case errors.Is(err, context.DeadlineExceeded):
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to delete a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error(err.Error(),
				zap.Stack("stackError"),
				zap.String("recipeId", recipeId.String()),
				zap.String("method", ctx.Request.Method),
				zap.String("path", ctx.FullPath()),
			)
		}
		return
	}

	h.logger.Info("deleted a recipe from database",
		zap.String("recipeId", recipeId.String()),
	)

	ctx.Status(http.StatusNoContent)
}

func (h *handler) getRecipeById(ctx *gin.Context) {
	recipeIdParam, ok := ctx.Params.Get("id")

	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "missing ID param for recipe",
		})
		return
	}

	recipeId, err := uuid.Parse(recipeIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "the received ID is not a valid UUID",
		})
		return
	}

	r, err := h.recipeService.GetRecipeById(ctx.Copy(), recipeId)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "could not find recipe with this UUID",
			})
		case errors.Is(err, context.DeadlineExceeded):
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to find a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error(err.Error(),
				zap.Stack("stackError"),
				zap.String("recipeId", recipeId.String()),
				zap.String("method", ctx.Request.Method),
				zap.String("path", ctx.FullPath()),
			)
		}
		return
	}

	dto := RecipeResponse{
		Title:     r.Title,
		Content:   r.Content,
		CreatedAt: r.CreatedAt.Time,
		UpdatedAt: r.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, shared.DataResponse[RecipeResponse]{
		Data: dto,
	})
}
