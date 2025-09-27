package user

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
)

func (e *EntUserClient) FindByID(ctx context.Context, id int) (*ent.User, error) {
	user, err := e.client.User().
		Query().
		Where(user.ID(id)).
		Where(user.DeletedAtIsNil()).
		First(ctx)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
