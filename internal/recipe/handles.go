package recipe

import (
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

	if validateRecipe(v, requestBody); !v.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request did not pass the validation",
			"fields":  v.Errors,
		})
		return
	}

	recipeId, err := h.recipeService.createNewRecipe(ctx.Copy(), requestBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong when saving a recipe",
		})
		return
	}

	// TODO: find out a better way to get the address
	ctx.Header("Location", fmt.Sprintf("http://localhost:8080/api/v1/recipes/%v", recipeId.String()))
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "successfully saved a recipe",
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
			"message": "the recieved ID is not a valid UUID",
		})
		return
	}

	_ = h.recipeService.deleteRecipeById(ctx.Copy(), recipeId)

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
			"message": "the recieved ID is not a valid UUID",
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
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "failed to return recipe",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
		"recipe": r,
		},
	})
}
