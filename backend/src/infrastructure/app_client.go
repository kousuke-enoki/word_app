// infrastructure/app_client.go
package infrastructure

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
	"word_app/backend/src/interfaces"
	"word_app/backend/src/models"
)

type appClient struct {
	entClient *ent.Client
}

// DeleteWord implements interfaces.ClientInterface.
func (c *appClient) DeleteWord(ctx context.Context, userID int, wordID int) (*models.DeleteWordResponse, error) {
	panic("unimplemented")
}

// GetRegisteredWords implements interfaces.ClientInterface.
func (c *appClient) GetRegisteredWords(ctx context.Context, AllWordListRequest *models.AllWordListRequest) (*models.AllWordListResponse, error) {
	panic("unimplemented")
}

// GetWordDetails implements interfaces.ClientInterface.
func (c *appClient) GetWordDetails(ctx context.Context, wordID int, userID int) (*models.WordShowResponse, error) {
	panic("unimplemented")
}

// GetWords implements interfaces.ClientInterface.
func (c *appClient) GetWords(ctx context.Context, AllWordListRequest *models.AllWordListRequest) (*models.AllWordListResponse, error) {
	panic("unimplemented")
}

// RegisterWords implements interfaces.ClientInterface.
func (c *appClient) RegisterWords(ctx context.Context, wordID int, userID int, IsRegistered bool) (*models.RegisterWordResponse, error) {
	panic("unimplemented")
}

// SaveMemo implements interfaces.ClientInterface.
func (c *appClient) SaveMemo(ctx context.Context, wordID int, userID int, memo string) (*models.SaveMemoResponse, error) {
	panic("unimplemented")
}

// UpdateWord implements interfaces.ClientInterface.
func (c *appClient) UpdateWord(ctx context.Context, UpdateWordRequest *models.UpdateWordRequest) (*models.UpdateWordResponse, error) {
	panic("unimplemented")
}

// NewAppClient 初期化関数
func NewAppClient(entClient *ent.Client) interfaces.ClientInterface {
	return &appClient{
		entClient: entClient,
	}
}

// UserClient の実装
func (c *appClient) CreateUser(ctx context.Context, email, name, password string) (*ent.User, error) {
	return c.entClient.User.Create().
		SetEmail(email).
		SetName(name).
		SetPassword(password).
		Save(ctx)
}

func (c *appClient) FindUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	return c.entClient.User.Query().Where(user.Email(email)).Only(ctx)
}

func (c *appClient) FindUserByID(ctx context.Context, id int) (*ent.User, error) {
	return c.entClient.User.Query().Where(user.ID(id)).Only(ctx)
}

// WordService の実装
func (c *appClient) CreateWord(ctx context.Context, req *models.CreateWordRequest) (*models.CreateWordResponse, error) {
	// 実装例
	return nil, nil
}

// EntClient を返す
func (c *appClient) EntClient() *ent.Client {
	return c.entClient
}
