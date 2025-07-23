package quiz

import (
	"errors"

	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

type ServiceImpl struct {
	client serviceinterfaces.EntClientInterface
}

func NewService(client serviceinterfaces.EntClientInterface) *ServiceImpl {
	return &ServiceImpl{client: client}
}

var (
	ErrQuizNotFound    = errors.New("word not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrDeleteQuiz      = errors.New("failed to delete word")
	ErrDatabaseFailure = errors.New("database failure")
)
