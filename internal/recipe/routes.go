package recipe

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) RegisterRoutes(gin *gin.Engine) {

	gin.POST("api/v1/recipes", h.createRecipe)
	gin.GET("api/v1/recipes/:id", h.getRecipeById)
	gin.DELETE("api/v1/recipes/:id", h.deleteRecipeById)
}
