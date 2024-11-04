package models

import (
	"context"
	"word_app/backend/ent"
)

type Client struct {
	*ent.Client
}

func (c *Client) User() *ent.UserClient {
	return c.Client.User
}

func (c *Client) UserCreate(ctx context.Context, email string, name string, password string) (*ent.User, error) {
	return c.User().
		Create().
		SetEmail(email).
		SetName(name).
		SetPassword(password).
		Save(ctx)
}
