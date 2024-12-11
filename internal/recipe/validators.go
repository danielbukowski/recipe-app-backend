package recipe

import "github.com/danielbukowski/recipe-app-backend/internal/validator"

type newRecipeRequest struct {
	title   string
	content string
}

func validateRecipe(v *validator.Validator, requestBody newRecipeRequest) {
	v.Check(len(requestBody.content) <= 10, "content", "should be at least 10 characters length long")
	v.Check(len(requestBody.title) <= 10, "title", "should be at least 10 characters length long")
}
