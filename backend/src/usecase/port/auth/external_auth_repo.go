package auth

import (
	"context"
	"word_app/backend/src/domain"
)

type ExternalAuthRepository interface {
	Create(ctx context.Context, ext *domain.ExternalAuth) error
}
