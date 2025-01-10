package word_service

import (
	"errors"
	"word_app/backend/ent"
)

type WordServiceImpl struct {
	client *ent.Client
}

func NewWordService(client *ent.Client) *WordServiceImpl {
	return &WordServiceImpl{client: client}
}

var (
	ErrWordNotFound    = errors.New("word not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrDeleteWord      = errors.New("failed to delete word")
	ErrDatabaseFailure = errors.New("database failure")
)
