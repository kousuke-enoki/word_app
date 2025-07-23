package quiz

import (
	"word_app/backend/src/interfaces/http/quiz"
)

type Handler struct {
	quizService quiz.Service
}

func NewHandler(quizService quiz.Service) *Handler {
	return &Handler{quizService: quizService}
}
