package user

import (
	"context"

	"word_app/backend/ent/externalauth"
	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
)

func (r *EntUserRepo) FindByProvider(ctx context.Context, provider, subject string) (*domain.User, error) {
	u, err := r.client.User().
		Query().
		Where(
			user.HasExternalAuthsWith(
				externalauth.Provider(provider),
				externalauth.ProviderUserID(subject),
			)).
		Only(ctx)
	if err != nil {
		return nil, err // ent.ErrNotFound なら呼び出し側で nil 判定
	}
	return &domain.User{ID: u.ID, Email: u.Email, Name: u.Name}, nil
}
