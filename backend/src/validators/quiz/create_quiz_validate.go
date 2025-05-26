package validator

import (
	"word_app/backend/src/models"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateCreateQuiz(dto *models.CreateQuizDTO) error {
	return v.Struct(dto)
}
