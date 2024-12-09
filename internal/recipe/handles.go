package recipe

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type controller struct {
	Logger        *zap.Logger
	recipeService *service
}

func NewController(logger *zap.Logger, recipeService *service) *controller {
	return &controller{
		Logger:        logger,
		recipeService: recipeService,
	}
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
