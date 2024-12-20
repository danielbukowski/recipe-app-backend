package recipe

import "time"

type RecipeResponse struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewRecipeRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateRecipeRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
