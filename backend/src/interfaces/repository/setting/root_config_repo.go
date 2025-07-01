package setting

import (
	"context"
	"word_app/backend/src/domain"
)

type RootConfigRepository interface {
	Get(ctx context.Context) (*domain.RootConfig, error)
	Upsert(ctx context.Context, cfg *domain.RootConfig) (*domain.RootConfig, error)
}
