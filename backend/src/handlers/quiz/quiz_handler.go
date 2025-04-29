package quiz

import (
	"word_app/backend/src/interfaces"
)

type QuizHandler struct {
	quizService interfaces.QuizService
}

func NewQuizHandler(quizService interfaces.QuizService) *QuizHandler {
	return &QuizHandler{quizService: quizService}
}
