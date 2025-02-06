package recipe

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes sets endpoints for Recipe resource.
func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.POST("api/v1/recipes", h.CreateRecipe)
	e.GET("api/v1/recipes/:id", h.GetRecipeById)
	e.PUT("api/v1/recipes/:id", h.UpdateRecipeById)
	e.DELETE("api/v1/recipes/:id", h.DeleteRecipeById)
}
