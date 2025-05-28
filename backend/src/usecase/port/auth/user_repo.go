package auth

import (
	"context"
	"word_app/backend/src/domain"
	"word_app/backend/src/models"
)

type UserRepository interface {
	FindByProvider(ctx context.Context, provider, sub string) (*models.User, error)
	Create(ctx context.Context, u *domain.User, ext *domain.ExternalAuth) error
}
