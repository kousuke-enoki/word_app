package handlers

import (
	"errors"
	"strings"

	"word_app/backend/src/models"

	"github.com/go-playground/validator/v10"
)

func FieldsFromBindError(err error) []models.FieldError {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return nil
	}
	out := make([]models.FieldError, 0, len(verrs))
	for _, fe := range verrs {
		field := fe.Field() // 例: "Email"
		// JSON名に寄せたいならタグから取る処理を追加してもOK
		msg := ""
		switch fe.Tag() {
		case "required":
			msg = "is required"
		case "email":
			msg = "invalid email format"
		default:
			msg = fe.Error()
		}
		out = append(out, models.FieldError{
			Field:   strings.ToLower(field), // "email"など
			Message: msg,
		})
	}
	return out
}
