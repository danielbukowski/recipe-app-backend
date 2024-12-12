package recipe

import "time"

type newRecipeRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
