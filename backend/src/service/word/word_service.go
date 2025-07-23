package word

import (
	"errors"

	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

type ServiceImpl struct {
	client serviceinterfaces.EntClientInterface
}

func NewWordService(client serviceinterfaces.EntClientInterface) *ServiceImpl {
	return &ServiceImpl{client: client}
}

var (
	ErrWordNotFound       = errors.New("word not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrDeleteWord         = errors.New("failed to delete word")
	ErrDatabaseFailure    = errors.New("database failure")
	ErrWordExists         = errors.New("there is already a word with the same name")
	ErrCreateWord         = errors.New("failed to create word")
	ErrCreateWordInfo     = errors.New("failed to create word info")
	ErrCreateJapaneseMean = errors.New("failed to create japanese mean")
)
