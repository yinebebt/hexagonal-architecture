package util

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func CoolTitleValidator(field validator.FieldLevel) bool {
	return strings.Contains(field.Field().String(), "cool")
}
