package user

import (
	"context"
	"word_app/backend/ent"
)

// UserClientInterface ユーザーエンティティ用のインターフェース
type UserClientInterface interface {
	User() *ent.UserClient // UserClientをポインタ型で返す
	Create() UserCreateInterface
}

// UserCreateInterface ユーザー作成用のインターフェース
type UserCreateInterface interface {
	SetEmail(string) UserCreateInterface
	SetName(string) UserCreateInterface
	SetPassword(string) UserCreateInterface
	Save(context.Context) (*ent.User, error)
}

// ClientInterface
type ClientInterface interface {
	User() *ent.UserClient
	UserCreate(ctx context.Context, email string, name string, password string) (*ent.User, error)
}
