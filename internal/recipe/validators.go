package recipe

import (
	"github.com/danielbukowski/recipe-app-backend/internal/validator"
)

type newRecipeRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func validateRecipe(v *validator.Validator, requestBody newRecipeRequest) {
	v.Check(len(requestBody.Content) >= 10, "content", "should be at least 10 characters length long")
	v.Check(len(requestBody.Title) >= 10, "title", "should be at least 10 characters length long")
}
