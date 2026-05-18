package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func FormatValidationError(err error) map[string]string {

	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {

		field := strings.ToLower(err.Field())

		switch err.Tag() {

		case "required":
			errors[field] = field + " is required"

		case "email":
			errors[field] = "invalid email format"

		case "min":
			errors[field] = field + " is too short"

		case "max":
			errors[field] = field + " is too long"

		default:
			errors[field] = "invalid value"
		}
	}

	return errors
}
