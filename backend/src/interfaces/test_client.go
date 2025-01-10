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
func (w *TestClientWrapper) GetRegisteredWords(ctx context.Context, AllWordListRequest *models.AllWordListRequest) (*models.AllWordListResponse, error) {
	panic("unimplemented")
}

// GetWordDetails implements ClientInterface.
func (w *TestClientWrapper) GetWordDetails(ctx context.Context, WordShowRequest *models.WordShowRequest) (*models.WordShowResponse, error) {
	panic("unimplemented")
}

// GetWords implements ClientInterface.
func (w *TestClientWrapper) GetWords(ctx context.Context, AllWordListRequest *models.AllWordListRequest) (*models.AllWordListResponse, error) {
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
