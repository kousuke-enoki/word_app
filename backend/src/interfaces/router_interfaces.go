package interfaces

import (
	"context"

	"word_app/backend/ent"
)

// インターフェース定義
type SampleUserClient interface {
	Query() UserQuery
	SchemaCreate(ctx context.Context) error
}

// Query メソッドで返す型
type UserQuery interface {
	Where(predicates ...func(*ent.User) bool) UserQuery
	Exist(ctx context.Context) (bool, error)
}

type UserQueryInterface interface {
	Where(predicate ...interface{}) UserQueryInterface
	Exist(ctx context.Context) (bool, error)
}
