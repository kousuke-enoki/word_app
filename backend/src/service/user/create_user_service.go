package user

import (
	"context"

	"word_app/backend/ent"
)

func (e *EntUserClient) Create(ctx context.Context, email, name, password string) (*ent.User, error) {
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

	_, err = e.client.UserConfig().
		Create().
		SetUserID(user.ID).
		SetIsDarkMode(false).
		Save(ctx)

	if err != nil {
		return nil, ErrCreateConfig
	}

	return user, nil
}
