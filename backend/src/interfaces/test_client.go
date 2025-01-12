package interfaces

import (
	"context"
	"word_app/backend/ent"
	"word_app/backend/src/models"
)

type TestClientWrapper struct {
	entClient *ent.Client
}

// CreateUser implements ClientInterface.
func (w *TestClientWrapper) CreateUser(ctx context.Context, email string, name string, password string) (*ent.User, error) {
	panic("unimplemented")
}

// CreateWord implements ClientInterface.
func (w *TestClientWrapper) CreateWord(ctx context.Context, CreateWordRequest *models.CreateWordRequest) (*models.CreateWordResponse, error) {
	panic("unimplemented")
}

// DeleteWord implements ClientInterface.
func (w *TestClientWrapper) DeleteWord(ctx context.Context, DeleteWordRequest *models.DeleteWordRequest) (*models.DeleteWordResponse, error) {
	panic("unimplemented")
}

// FindUserByEmail implements ClientInterface.
func (w *TestClientWrapper) FindUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	panic("unimplemented")
}

// FindUserByID implements ClientInterface.
func (w *TestClientWrapper) FindUserByID(ctx context.Context, id int) (*ent.User, error) {
	panic("unimplemented")
}

// GetRegisteredWords implements ClientInterface.
func (w *TestClientWrapper) GetRegisteredWords(ctx context.Context, WordListRequest *models.WordListRequest) (*models.WordListResponse, error) {
	panic("unimplemented")
}

// GetWordDetails implements ClientInterface.
func (w *TestClientWrapper) GetWordDetails(ctx context.Context, WordShowRequest *models.WordShowRequest) (*models.WordShowResponse, error) {
	panic("unimplemented")
}

// GetWords implements ClientInterface.
func (w *TestClientWrapper) GetWords(ctx context.Context, WordListRequest *models.WordListRequest) (*models.WordListResponse, error) {
	panic("unimplemented")
}

// RegisterWords implements ClientInterface.
func (w *TestClientWrapper) RegisterWords(ctx context.Context, RegisterWordRequest *models.RegisterWordRequest) (*models.RegisterWordResponse, error) {
	panic("unimplemented")
}

// SaveMemo implements ClientInterface.
func (w *TestClientWrapper) SaveMemo(ctx context.Context, SaveMemoRequest *models.SaveMemoRequest) (*models.SaveMemoResponse, error) {
	panic("unimplemented")
}

// UpdateWord implements ClientInterface.
func (w *TestClientWrapper) UpdateWord(ctx context.Context, UpdateWordRequest *models.UpdateWordRequest) (*models.UpdateWordResponse, error) {
	panic("unimplemented")
}

func NewTestClientWrapper(client *ent.Client) ClientInterface {
	return &TestClientWrapper{entClient: client}
}

func (w *TestClientWrapper) EntClient() *ent.Client {
	return w.entClient
}

// Tx はトランザクションを開始します。
func (w *TestClientWrapper) Tx(ctx context.Context) (*ent.Tx, error) {
	panic("unimplemented")
}

// Word は WordClient を返します。
func (w *TestClientWrapper) Word() *ent.WordClient {
	panic("unimplemented")
}

// User は UserClient を返します。
func (w *TestClientWrapper) User() *ent.UserClient {
	panic("unimplemented")
}

// RegisteredWord は RegisteredWordClient を返します。
func (w *TestClientWrapper) RegisteredWord() *ent.RegisteredWordClient {
	panic("unimplemented")
}

// WordInfo は WordInfoClient を返します。
func (w *TestClientWrapper) WordInfo() *ent.WordInfoClient {
	panic("unimplemented")
}

// JapaneseMean は JapaneseMeanClient を返します。
func (w *TestClientWrapper) JapaneseMean() *ent.JapaneseMeanClient {
	panic("unimplemented")
}
