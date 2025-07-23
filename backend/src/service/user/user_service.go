package user

import (
	"errors"

	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

// ent.Client をラップして UserClient インターフェースを実装
type EntUserClient struct {
	client serviceinterfaces.EntClientInterface
}

func NewEntUserClient(client serviceinterfaces.EntClientInterface) *EntUserClient {
	return &EntUserClient{client: client}
}

var (
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrDuplicateID     = errors.New("duplicate ID")
	ErrDatabaseFailure = errors.New("database failure")
	ErrUserNotFound    = errors.New("user not found")
	ErrCreateConfig    = errors.New("create user config failure")
)
