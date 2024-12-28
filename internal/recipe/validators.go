package recipe

import (
	"fmt"
	"unicode/utf8"

	"github.com/danielbukowski/recipe-app-backend/internal/validator"
)

const (
	minTitleLength   = 4
	minContentLength = 10
)

func validateNewRecipeRequestBody(v *validator.Validator, requestBody NewRecipeRequest) {
	v.Check(utf8.RuneCountInString(requestBody.Content) >= minContentLength, "content", fmt.Sprintf("should be at least %d characters length long", minContentLength))
	v.Check(utf8.RuneCountInString(requestBody.Title) >= minTitleLength, "title", fmt.Sprintf("should be at least %d characters length long", minTitleLength))
}

func validateUpdateRecipeRequestBody(v *validator.Validator, requestBody UpdateRecipeRequest) {
	v.Check(utf8.RuneCountInString(requestBody.Content) >= minContentLength, "content", fmt.Sprintf("should be at least %d characters length long", minContentLength))
	v.Check(utf8.RuneCountInString(requestBody.Title) >= minTitleLength, "title", fmt.Sprintf("should be at least %d characters length long", minTitleLength))
}
