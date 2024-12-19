package recipe

import (
	"fmt"

	"github.com/danielbukowski/recipe-app-backend/internal/validator"
)

const (
	minTitleLength   = 4
	minContentLength = 10
)

func validateNewRecipeRequestBody(v *validator.Validator, requestBody newRecipeRequest) {
	v.Check(len(requestBody.Content) >= minContentLength, "content", fmt.Sprintf("should be at least %d characters length long", minContentLength))
	v.Check(len(requestBody.Title) >= minTitleLength, "title", fmt.Sprintf("should be at least %d characters length long", minTitleLength))
}
