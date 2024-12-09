package recipe

import (
	"github.com/gin-gonic/gin"
)

func (c *controller) RegisterRoutes(gin *gin.Engine) {
	gin.GET("api/v1/recipes/:id", c.getRecipeByIdHandler)
}
