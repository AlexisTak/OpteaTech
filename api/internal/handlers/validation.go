package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func formatValidationErrors(err error) map[string]string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{"error": "validation failed"}
	}

	formatted := make(map[string]string, len(validationErrors))
	for _, item := range validationErrors {
		formatted[item.Field()] = fmt.Sprintf("validation failed on %s", item.Tag())
	}

	return formatted
}
