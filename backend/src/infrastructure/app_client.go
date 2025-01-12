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

// NewAppClient 初期化関数
func NewAppClient(entClient *ent.Client) interfaces.ClientInterface {
	return &appClient{
		entClient: entClient,
	}
}

// DeleteWord implements interfaces.ClientInterface.
func (c *appClient) DeleteWord(ctx context.Context, DeleteWordRequest *models.DeleteWordRequest) (*models.DeleteWordResponse, error) {
	panic("unimplemented")
}

// GetRegisteredWords implements interfaces.ClientInterface.
func (c *appClient) GetRegisteredWords(ctx context.Context, WordListRequest *models.WordListRequest) (*models.WordListResponse, error) {
	panic("unimplemented")
}

// GetWordDetails implements interfaces.ClientInterface.
func (c *appClient) GetWordDetails(ctx context.Context, WordShowRequest *models.WordShowRequest) (*models.WordShowResponse, error) {
	panic("unimplemented")
}

// GetWords implements interfaces.ClientInterface.
func (c *appClient) GetWords(ctx context.Context, WordListRequest *models.WordListRequest) (*models.WordListResponse, error) {
	panic("unimplemented")
}

// RegisterWords implements interfaces.ClientInterface.
func (c *appClient) RegisterWords(ctx context.Context, RegisterWordRequest *models.RegisterWordRequest) (*models.RegisterWordResponse, error) {
	panic("unimplemented")
}

// SaveMemo implements interfaces.ClientInterface.
func (c *appClient) SaveMemo(ctx context.Context, SaveMemoRequest *models.SaveMemoRequest) (*models.SaveMemoResponse, error) {
	panic("unimplemented")
}

// UpdateWord implements interfaces.ClientInterface.
func (c *appClient) UpdateWord(ctx context.Context, UpdateWordRequest *models.UpdateWordRequest) (*models.UpdateWordResponse, error) {
	panic("unimplemented")
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

// Tx はトランザクションを開始します。
func (c *appClient) Tx(ctx context.Context) (*ent.Tx, error) {
	return c.entClient.Tx(ctx)
}

// Word は WordClient を返します。
func (c *appClient) Word() *ent.WordClient {
	return c.entClient.Word
}

// User は UserClient を返します。
func (c *appClient) User() *ent.UserClient {
	return c.entClient.User
}

// RegisteredWord は RegisteredWordClient を返します。
func (c *appClient) RegisteredWord() *ent.RegisteredWordClient {
	return c.entClient.RegisteredWord
}

// WordInfo は WordInfoClient を返します。
func (c *appClient) WordInfo() *ent.WordInfoClient {
	return c.entClient.WordInfo
}

// JapaneseMean は JapaneseMeanClient を返します。
func (c *appClient) JapaneseMean() *ent.JapaneseMeanClient {
	return c.entClient.JapaneseMean
}
