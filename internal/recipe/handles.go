package recipe

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danielbukowski/recipe-app-backend/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type handler struct {
	logger        *zap.Logger
	recipeService *service
}

func NewHandler(logger *zap.Logger, recipeService *service) *handler {
	return &handler{
		logger:        logger,
		recipeService: recipeService,
	}
}

func (h *handler) createRecipe(ctx *gin.Context) {
	var requestBody = newRecipeRequest{}

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse the request body",
		})
		return
	}

	v := validator.New()

	if validateNewRecipeRequestBody(v, requestBody); !v.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request did not pass the validation",
			"fields":  v.Errors,
		})
		return
	}

	recipeId, err := h.recipeService.createNewRecipe(ctx.Copy(), requestBody)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to save a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error("createNewRecipe method threw unexpected behavior",
				zap.String("recipeId", recipeId.String()),
			)
		}
		return
	}

	// TODO: find out a better way to get the address
	ctx.Header("Location", fmt.Sprintf("http://localhost:8080/api/v1/recipes/%v", recipeId.String()))
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "successfully saved a recipe",
	})
}

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

	var requestBody = newRecipeRequest{}

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse the request body",
		})
		return
	}

	recipe, err := h.recipeService.getRecipeById(ctx, recipeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "could not find a recipe with this UUID",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	if requestBody.Title == "" {
		requestBody.Title = recipe.Title
	}

	if requestBody.Content == "" {
		requestBody.Content = recipe.Content
	}

	v := validator.New()
	validateNewRecipeRequestBody(v, requestBody)

	if validateNewRecipeRequestBody(v, requestBody); !v.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request did not pass the validation",
			"fields":  v.Errors,
		})
		return
	}

	err = h.recipeService.updateRecipeById(ctx.Copy(), recipeId, recipe.UpdatedAt, requestBody)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			ctx.JSON(http.StatusConflict, gin.H{
				"message": "conflict occurred when trying to update a recipe",
			})
		case context.DeadlineExceeded:
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to save a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error("updateRecipeById method threw unexpected behavior",
				zap.String("recipeId", recipeId.String()),
			)
		}
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "successfully updated a recipe",
	})

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

	err = h.recipeService.deleteRecipeById(ctx.Copy(), recipeId)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			ctx.Status(http.StatusNoContent)
		case context.DeadlineExceeded:
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to delete a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error("deleteRecipeById method threw unexpected behavior",
				zap.String("recipeId", recipeId.String()),
			)
		}
		return
	}

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

	r, err := h.recipeService.getRecipeById(ctx.Copy(), recipeId)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "could not find recipe with this UUID",
			})
		case context.DeadlineExceeded:
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to fetch a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error("getRecipeById method threw unexpected behavior",
				zap.String("recipeId", recipeId.String()),
			)
		}
		return
	}

	dto := recipeResponse{
		Title:     r.Title,
		Content:   r.Content,
		CreatedAt: r.CreatedAt.Time,
		UpdatedAt: r.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"recipe": dto,
		},
	})
}
