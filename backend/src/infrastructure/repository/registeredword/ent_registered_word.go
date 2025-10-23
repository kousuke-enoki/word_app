// src/infrastructure/repository/registeredword/ent_registered_word.go
package registeredword

import (
	"context"

	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

type EntRegisteredWordReadRepo struct {
	client serviceinterfaces.EntClientInterface
}

func NewEntRegisteredWordReadRepo(
	client serviceinterfaces.EntClientInterface,
) *EntRegisteredWordReadRepo {
	return &EntRegisteredWordReadRepo{
		client: client,
	}
}

type EntRegisteredWordWriteRepo struct {
	client serviceinterfaces.EntClientInterface
}

func NewEntRegisteredWordWriteRepo(
	client serviceinterfaces.EntClientInterface,
) *EntRegisteredWordWriteRepo {
	return &EntRegisteredWordWriteRepo{
		client: client,
	}
}

// RegisteredWord リポジトリ
type ReadRepository interface {
	// 指定 userID の「active=true」な registered_word の WordID セットを返す
	ActiveWordIDSetByUser(ctx context.Context, userID int, wordIDs []int) (map[int]struct{}, error)

	// 現在 active=true の登録総数
	CountActiveByUser(ctx context.Context, userID int) (int, error)
	// 指定 userID の registered_word（active/非active問わず）の存在有無を返す
	// map[wordID]isActive
	FindActiveMapByUserAndWordIDs(ctx context.Context, userID int, wordIDs []int) (map[int]bool, error)
}

// RegisteredWord リポジトリ
type WriteRepository interface {
	// 既存→ is_active=true に更新。存在しなければ ErrNotFound
	Activate(ctx context.Context, userID, wordID int) error
	// 新規作成。UNIQUE(user_id, word_id) 競合時は ErrConflict を返す実装でもOK
	CreateActive(ctx context.Context, userID, wordID int) error
}
