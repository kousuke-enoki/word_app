package result

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
	ErrResultNotFound  = errors.New("word not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrDeleteResult    = errors.New("failed to delete word")
	ErrDatabaseFailure = errors.New("database failure")
)
