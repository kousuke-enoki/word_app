package user_service

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
)

func (e *EntUserClient) FindUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	user, err := e.client.User().
		Query().
		Where(user.EmailEQ(email)).
		First(ctx)

	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
