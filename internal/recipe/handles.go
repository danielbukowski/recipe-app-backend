package recipe

import (
	"context"
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

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
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

			h.logger.Error("createNewRecipe method in createRecipe handler threw unexpected behavior",
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

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "missing JSON request body",
		})
		return
	}

	recipeFromDb, err := h.recipeService.getRecipeById(ctx, recipeId)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "could not find a recipe with this id",
			})
		case context.DeadlineExceeded:
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to update a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error("getRecipeById method in updateRecipeById handler threw unexpected behavior",
				zap.String("recipeId", recipeId.String()),
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
	validateNewRecipeRequestBody(v, requestBody)

	if validateNewRecipeRequestBody(v, requestBody); !v.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request body did not pass the validation",
			"fields":  v.Errors,
		})
		return
	}

	err = h.recipeService.updateRecipeById(ctx.Copy(), recipeId, recipeFromDb.UpdatedAt, requestBody)
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

			h.logger.Error("updateRecipeById method in updateRecipeById threw unexpected behavior",
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

			h.logger.Error("deleteRecipeById method in deleteRecipeById handler threw unexpected behavior",
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
				"message": "failed to find a recipe in time",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error("getRecipeById method in getRecipeById handler threw unexpected behavior",
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
