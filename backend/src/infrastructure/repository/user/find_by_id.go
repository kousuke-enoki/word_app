package user

import (
	"context"

	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
)

func (r *EntUserRepo) FindByID(ctx context.Context, id int) (*domain.User, error) {
	u, err := r.client.User().
		Query().
		Where(user.ID(id)).
		Select(user.FieldID, user.FieldIsRoot).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.User{ID: u.ID, IsRoot: u.IsRoot}, nil
}
