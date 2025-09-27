package user

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/externalauth"
	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
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
		return nil, ErrUserNotFound
	}
	return mapEntToDomain(u, u.Edges.ExternalAuths), nil
}

// --- mapper（Ent → Domain） ---
func mapEntToDomain(u *ent.User, auths []*ent.ExternalAuth) *domain.User {
	emailPtr := u.Email       // そのまま使える
	hasLine := len(auths) > 0 // QF1007対応

	return &domain.User{
		ID:          u.ID,
		Email:       emailPtr,
		Name:        u.Name,
		IsAdmin:     u.IsAdmin,
		IsRoot:      u.IsRoot,
		IsTest:      u.IsTest,
		HasPassword: u.Password != nil && *u.Password != "",
		HasLine:     hasLine,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
