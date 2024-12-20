package recipe

import "time"

type recipeResponse struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type newRecipeRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateRecipeRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
