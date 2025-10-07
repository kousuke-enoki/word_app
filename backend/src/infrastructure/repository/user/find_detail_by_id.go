package user

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/externalauth"
	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
	usermapper "word_app/backend/src/infrastructure/mapper/user"
	"word_app/backend/src/infrastructure/repoerr"
)

func (e *EntUserRepo) FindDetailByID(ctx context.Context, id int) (*domain.User, error) {
	u, err := e.client.User().
		Query().
		Where(user.ID(id), user.DeletedAtIsNil()).
		WithExternalAuths(func(q *ent.ExternalAuthQuery) {
			q.Where(externalauth.ProviderEqualFold("line"))
		}).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, repoerr.FromEnt(err, "user not found", "duplicate id") // 外部認証未連携
		}
		return nil, repoerr.FromEnt(err, "internal", "internal server error")
	}
	return usermapper.MapEntUser(u, usermapper.WithAuths(u.Edges.ExternalAuths)), nil
}
