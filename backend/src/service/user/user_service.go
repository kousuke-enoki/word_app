package user_service

import (
	"errors"
	"word_app/backend/src/interfaces/service_interfaces"
)

// ent.Client をラップして UserClient インターフェースを実装
type EntUserClient struct {
	client service_interfaces.EntClientInterface
}

func NewEntUserClient(client service_interfaces.EntClientInterface) *EntUserClient {
	return &EntUserClient{client: client}
}

var (
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrDuplicateID     = errors.New("duplicate ID")
	ErrDatabaseFailure = errors.New("database failure")
	ErrUserNotFound    = errors.New("user not found")
)
