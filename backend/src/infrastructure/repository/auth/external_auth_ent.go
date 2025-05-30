// infrastructure/repository/external_auth_ent.go
package auth

import (
	"context"

	"word_app/backend/src/domain"
	"word_app/backend/src/interfaces/service_interfaces"
)

type EntExtAuthRepo struct {
	client service_interfaces.EntClientInterface
}

func NewEntExtAuthRepo(c service_interfaces.EntClientInterface) *EntExtAuthRepo {
	return &EntExtAuthRepo{client: c}
}

func (r *EntExtAuthRepo) Create(ctx context.Context, ext *domain.ExternalAuth) error {
	_, err := r.client.ExternalAuth().
		Create().
		SetUserID(ext.UserID).
		SetProvider(ext.Provider).
		SetProviderUserID(ext.ProviderUserID).
		Save(ctx)
	return err
}
