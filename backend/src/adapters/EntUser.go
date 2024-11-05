package adapters

import (
	"context"
	"word_app/backend/ent"
	"word_app/backend/ent/user"
)

// ent.Client をラップして UserClient インターフェースを実装
type EntUserClient struct {
	client *ent.Client
}

func NewEntUserClient(client *ent.Client) *EntUserClient {
	return &EntUserClient{client: client}
}

func (e *EntUserClient) CreateUser(ctx context.Context, email, name, password string) (*ent.User, error) {
	return e.client.User.
		Create().
		SetEmail(email).
		SetName(name).
		SetPassword(password).
		Save(ctx)
}

func (e *EntUserClient) FindUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	return e.client.User.
		Query().
		Where(user.EmailEQ(email)).
		First(ctx)
}

func (e *EntUserClient) FindUserByID(ctx context.Context, id int) (*ent.User, error) { // 追加
	return e.client.User.
		Query().
		Where(user.ID(id)).
		First(ctx)
}
