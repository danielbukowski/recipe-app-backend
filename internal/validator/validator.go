package validator

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type validationError struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	return &Validator{validator: validator.New()}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {

		var vErr *validationError = &validationError{
			Message: "Your request body did not pass the validation",
			Fields:  make(map[string]string),
		}

		for _, err := range err.(validator.ValidationErrors) {
			var message string

			switch err.ActualTag() {
			case "required":
				message = "cannot be blank"
			case "email":
				message = "must be a valid email address"
			case "eqfield":
				message = fmt.Sprintf("must be the same as the %s field", err.Param())
			case "min":
				message = fmt.Sprintf("must be at least %v characters long", err.Param())
			case "max":
				message = fmt.Sprintf("cannot exceed %v characters", err.Param())
			default:
				message = fmt.Sprintf("Field '%s': '%v' must satisfy '%s' '%v' criteria", err.Field(), err.Value(), err.Tag(), err.Param())
			}

			vErr.Fields[err.Field()] = message
		}

		return echo.NewHTTPError(http.StatusBadRequest, vErr)
	}
	return nil
}
