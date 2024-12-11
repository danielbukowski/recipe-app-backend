package recipe

import (
	"net/http"

	"github.com/danielbukowski/recipe-app-backend/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type controller struct {
	logger        *zap.Logger
	recipeService *service
}

func NewController(logger *zap.Logger, recipeService *service) *controller {
	return &controller{
		logger:        logger,
		recipeService: recipeService,
	}
}

func (c *controller) createRecipe(ctx *gin.Context) {
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

	if err := c.recipeService.CreateNewRecipe(ctx.Copy(), requestBody); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong when saving a recipe",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully saved a recipe",
	})
}

func (c *controller) deleteRecipeById(ctx *gin.Context) {
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

	_ = c.recipeService.deleteRecipeById(ctx.Copy(), recipeId)

	ctx.Status(http.StatusNoContent)
}

func (c *controller) getRecipeByIdHandler(ctx *gin.Context) {
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

	r, err := c.recipeService.GetRecipeById(ctx.Copy(), recipeId)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			ctx.JSON(http.StatusNotFound, "could not find recipe with this UUID")
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "failed to return recipe",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"recipe": r,
	})
}
