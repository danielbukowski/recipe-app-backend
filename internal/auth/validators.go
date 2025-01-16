package auth

import (
	"fmt"
	"net/mail"

	"github.com/danielbukowski/recipe-app-backend/internal/validator"
)

const minPasswordLength = 5

func validateSignUpRequestBody(v *validator.Validator, requestBody SignUpRequest) {
	v.Check(len(requestBody.Password) >= 6, "password", fmt.Sprintf("should be at least %d characters long", minPasswordLength))

	v.Check(requestBody.Password == requestBody.PasswordAgain, "passwords", "should be the same")

	v.Check(len(requestBody.Email) > 0, "email", "email should not be empty")
	v.Check(isEmailValid(requestBody.Email), "email", "should have a correct format")
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
