package user_service

import (
	"context"
	"word_app/backend/ent"
	"word_app/backend/ent/user"
)

func (e *EntUserClient) FindUserByID(ctx context.Context, id int) (*ent.User, error) {
	user, err := e.client.User.
		Query().
		Where(user.ID(id)).
		First(ctx)

	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
