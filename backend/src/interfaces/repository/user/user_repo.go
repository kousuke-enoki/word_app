package auth

import (
	"context"
	"word_app/backend/src/domain"
)

type UserRepository interface {
	FindByProvider(ctx context.Context, provider, sub string) (*domain.User, error)
	Create(ctx context.Context, u *domain.User, ext *domain.ExternalAuth) error
	FindByID(ctx context.Context, id int) (*domain.User, error)
	IsRoot(ctx context.Context, userID int) (bool, error)
}
