package adapters

import (
	"context"
	"word_app/backend/ent"
)

type EntUserClient struct {
	client *ent.Client
}

func (e *EntUserClient) CreateUser(ctx context.Context, email, name, password string) (*ent.User, error) {
	return e.client.User.
		Create().
		SetEmail(email).
		SetName(name).
		SetPassword(password).
		Save(ctx)
}
