package validator

import (
	"errors"
	"fmt"

	playground "github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	var validationErrors playground.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return map[string]string{"body": err.Error()}
	}

	formatted := make(map[string]string, len(validationErrors))
	for _, fieldErr := range validationErrors {
		field := fieldErr.Field()
		formatted[field] = messageFor(fieldErr)
	}
	return formatted
}

func messageFor(err playground.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}
