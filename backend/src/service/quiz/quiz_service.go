package quiz_service

import (
	"errors"
	"word_app/backend/src/interfaces/service_interfaces"
)

type QuizServiceImpl struct {
	client service_interfaces.EntClientInterface
}

func NewQuizService(client service_interfaces.EntClientInterface) *QuizServiceImpl {
	return &QuizServiceImpl{client: client}
}

var (
	ErrQuizNotFound    = errors.New("word not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrDeleteQuiz      = errors.New("failed to delete word")
	ErrDatabaseFailure = errors.New("database failure")
)
