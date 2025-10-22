package word

import (
	"errors"

	"word_app/backend/src/infrastructure/repository/user"
	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

type ServiceImpl struct {
	client   serviceinterfaces.EntClientInterface
	userRepo user.Repository
}

func NewWordService(
	client serviceinterfaces.EntClientInterface,
	userRepo user.Repository,
) *ServiceImpl {
	return &ServiceImpl{
		client:   client,
		userRepo: userRepo,
	}
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
