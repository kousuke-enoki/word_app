package user

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
	"word_app/backend/src/infrastructure/repoerr"
)

func (r *EntUserRepo) FindByID(ctx context.Context, id int) (*domain.User, error) {
	u, err := r.client.User().
		Query().
		Where(user.ID(id)).
		Select(
			user.FieldID,
			user.FieldIsRoot,
			user.FieldIsAdmin,
			user.FieldIsTest,
			user.FieldDeletedAt,
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, repoerr.FromEnt(err, "user not found", "duplicate id")
		}
		return nil, repoerr.FromEnt(err, "internal", "internal server error")
	}

	return &domain.User{
		ID:        u.ID,
		IsRoot:    u.IsRoot,
		IsAdmin:   u.IsAdmin,
		IsTest:    u.IsTest,
		DeletedAt: u.DeletedAt,
	}, nil
}
