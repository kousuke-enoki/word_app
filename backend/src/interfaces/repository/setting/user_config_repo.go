package setting

import (
	"context"

	"word_app/backend/src/domain"
)

type UserConfigRepository interface {
	FindByUserID(ctx context.Context, id int) (*domain.UserConfig, error)
	GetByUserID(ctx context.Context, userID int) (*domain.UserConfig, error)
	Upsert(ctx context.Context, cfg *domain.UserConfig) (*domain.UserConfig, error)
}
