package result_service

import (
	"errors"

	"word_app/backend/src/interfaces/service_interfaces"
)

type ResultServiceImpl struct {
	client service_interfaces.EntClientInterface
}

func NewResultService(client service_interfaces.EntClientInterface) *ResultServiceImpl {
	return &ResultServiceImpl{client: client}
}

var (
	ErrResultNotFound  = errors.New("word not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrDeleteResult    = errors.New("failed to delete word")
	ErrDatabaseFailure = errors.New("database failure")
)
