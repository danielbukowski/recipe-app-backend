package recipe

import (
	"github.com/gin-gonic/gin"
)

func (c *controller) RegisterRoutes(gin *gin.Engine) {

	gin.POST("api/v1/recipes", c.createRecipe)
	gin.GET("api/v1/recipes/:id", c.getRecipeByIdHandler)
	gin.DELETE("api/v1/recipes/:id", c.deleteRecipeById)
}
