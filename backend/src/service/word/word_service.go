package word_service

import (
	"errors"

	"word_app/backend/src/interfaces/service_interfaces"
)

type WordServiceImpl struct {
	client service_interfaces.EntClientInterface
}

func NewWordService(client service_interfaces.EntClientInterface) *WordServiceImpl {
	return &WordServiceImpl{client: client}
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
