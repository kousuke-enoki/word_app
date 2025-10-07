package user

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/externalauth"
	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
	"word_app/backend/src/infrastructure/repoerr"
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
		if ent.IsNotFound(err) {
			return nil, repoerr.FromEnt(err, "external auth not found", "not linked") // 外部認証未連携
		}
		return nil, repoerr.FromEnt(err, "internal", "internal server error")
	}
	return &domain.User{ID: u.ID, Email: u.Email, Name: u.Name}, nil
}
