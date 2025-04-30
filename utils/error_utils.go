package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"motorcycleApp/domain/dto"
	"net/http"
)

func ParseValidationErrors(err error) dto.Error {
	var fieldErrors []dto.FieldError
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, ve := range validationErrors {
			fieldErrors = append(fieldErrors, dto.FieldError{
				Field:   ve.Field(),
				Message: validationMessage(ve),
			})
		}
	}

	return dto.Error{
		Code:        http.StatusBadRequest,
		Message:     "Validation Error",
		FieldErrors: fieldErrors,
	}
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Поле обязательно для заполнения"
	case "email":
		return "Неверный формат email"
	case "gte":
		return "Минимальная длина " + fe.Param()
	case "len":
		return "Должно быть длиной " + fe.Param()
	case "numeric":
		return "Должно быть числом"
	default:
		return "Неверное значение"
	}
}
