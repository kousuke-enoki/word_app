package user_service

import (
	"errors"

	"word_app/backend/ent"
)

// ent.Client をラップして UserClient インターフェースを実装
type EntUserClient struct {
	client *ent.Client
}

func NewEntUserClient(client *ent.Client) *EntUserClient {
	return &EntUserClient{client: client}
}

var (
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrDatabaseFailure = errors.New("database failure")
)
