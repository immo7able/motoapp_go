package utils

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"motorcycleApp/domain/dto"
)

func ParseValidationErrors(err error) dto.Error {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return dto.Error{
			Code:    http.StatusBadRequest,
			Message: "Ошибка валидации",
		}
	}

	fieldErrors := make([]dto.FieldError, 0, len(validationErrs))
	for _, fe := range validationErrs {
		fieldErrors = append(fieldErrors, dto.FieldError{
			Field:   fe.Field(),
			Message: translateValidationMessage(fe),
		})
	}

	return dto.Error{
		Code:        http.StatusBadRequest,
		Message:     "Ошибка валидации",
		FieldErrors: fieldErrors,
	}
}

func translateValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Поле обязательно для заполнения"
	case "email":
		return "Неверный формат email"
	case "gte":
		return "Минимальная длина — " + fe.Param()
	case "len":
		return "Должно содержать " + fe.Param() + " символов"
	case "numeric":
		return "Должно быть числом"
	default:
		return "Недопустимое значение"
	}
}
