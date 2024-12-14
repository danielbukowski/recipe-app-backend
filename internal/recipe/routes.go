package recipe

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) RegisterRoutes(gin *gin.Engine) {

	gin.POST("api/v1/recipes", h.createRecipeHandler)
	gin.GET("api/v1/recipes/:id", h.getRecipeByIdHandler)
	gin.DELETE("api/v1/recipes/:id", h.deleteRecipeByIdHandler)
}
