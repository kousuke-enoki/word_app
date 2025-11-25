// src/infrastructure/repository/word/ent_word.go
package word

import (
	"context"

	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

type EntWordReadRepo struct {
	client serviceinterfaces.EntClientInterface
}

func NewEntWordReadRepo(
	client serviceinterfaces.EntClientInterface,
) *EntWordReadRepo {
	return &EntWordReadRepo{client: client}
}

// Word リポジトリ
type ReadRepository interface {
	// names に含まれる単語の ID を map[name]id で返す
	FindIDsByNames(ctx context.Context, names []string) (map[string]int, error)
}
