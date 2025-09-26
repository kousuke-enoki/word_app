// infrastructure/repository/external_auth_ent.go
package auth

import (
	"context"
	"time"

	"word_app/backend/src/domain"
	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

type EntExtAuthRepo struct {
	client serviceinterfaces.EntClientInterface
}
type ExternalAuthRepository interface {
	Create(ctx context.Context, ext *domain.ExternalAuth) error
	SoftDeleteByUserID(ctx context.Context, userID int, t time.Time) error
}

func NewEntExtAuthRepo(c serviceinterfaces.EntClientInterface) *EntExtAuthRepo {
	return &EntExtAuthRepo{client: c}
}

func (r *EntExtAuthRepo) Create(ctx context.Context, ext *domain.ExternalAuth) error {
	_, err := r.client.ExternalAuth().
		Create().
		SetUserID(ext.UserID).
		SetProvider(ext.Provider).
		SetProviderUserID(ext.ProviderUserID).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}
