package user

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
)

func (e *EntUserClient) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	user, err := e.client.User().
		Query().
		Where(user.EmailEQ(email)).
		Where(user.DeletedAtIsNil()).
		First(ctx)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
