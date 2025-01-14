package user_service

import (
	"context"
	"word_app/backend/ent"
)

func (e *EntUserClient) CreateUser(ctx context.Context, email, name, password string) (*ent.User, error) {
	user, err := e.client.User().
		Create().
		SetEmail(email).
		SetName(name).
		SetPassword(password).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, ErrDuplicateEmail
		}
		return nil, ErrDatabaseFailure
	}
	return user, nil
}
