package recipe

import "time"

type RecipeResponse struct {
	Title     string    `json:"title" example:"Chocolate Cookies"`
	Content   string    `json:"content" example:"Having all your ingredients the same temperature really helps here"`
	CreatedAt time.Time `json:"created_at" example:"2025-02-05T21:35:31.00635Z"`
	UpdatedAt time.Time `json:"updated_at"  example:"2025-02-07T21:35:31.00635Z"`
}

type NewRecipeRequest struct {
	Title   string `json:"title" validate:"required,min=5" example:"Chocolate Cookies"`
	Content string `json:"content" validate:"required,min=5" example:"Having all your ingredients the same temperature really helps here"`
}

type UpdateRecipeRequest struct {
	Title   string `json:"title" validate:"required,min=5" example:"Chocolate Cookies"`
	Content string `json:"content" validate:"required,min=5" example:"Having all your ingredients the same temperature really helps here"`
}
