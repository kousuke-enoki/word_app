// infrastructure/app_client.go
package infrastructure

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
	"word_app/backend/src/interfaces"
	"word_app/backend/src/models"
	settingUc "word_app/backend/src/usecase/setting"
	"word_app/backend/src/utils/contextutil"
)

type appClient struct {
	entClient *ent.Client
}

// GetRootConfigExecute implements interfaces.ClientInterface.
func (w *appClient) GetRootConfigExecute(ctx context.Context, in settingUc.GetRootConfigInput) (*settingUc.GetRootConfigOutput, error) {
	panic("unimplemented")
}

// GetUserConfigExecute implements interfaces.ClientInterface.
func (w *appClient) GetUserConfigExecute(ctx context.Context, in settingUc.GetUserConfigInput) (*settingUc.GetUserConfigOutput, error) {
	panic("unimplemented")
}

// UpdateRootConfigExecute implements interfaces.ClientInterface.
func (w *appClient) UpdateRootConfigExecute(ctx context.Context, in settingUc.UpdateRootConfigInput) (*domain.RootConfig, error) {
	panic("unimplemented")
}

// UpdateUserConfigExecute implements interfaces.ClientInterface.
func (w *appClient) UpdateUserConfigExecute(ctx context.Context, in settingUc.UpdateUserConfigInput) (*domain.UserConfig, error) {
	panic("unimplemented")
}

// GetAuthConfig implements interfaces.ClientInterface.
func (w *appClient) GetAuthConfig(ctx context.Context) (*settingUc.AuthConfigDTO, error) {
	panic("unimplemented")
}

// GetAuthConfig implements interfaces.ClientInterface.
func (w *appClient) GetAuthConfigs(ctx context.Context) (*models.AuthSettingResponse, error) {
	panic("unimplemented")
}

// GetResultSummaries implements interfaces.ClientInterface.
func (w *appClient) Validate(ctx context.Context, tokenStr string) (contextutil.UserRoles, error) {
	panic("unimplemented")
}

// GetResultSummaries implements interfaces.ClientInterface.
func (w *appClient) GetResultSummaries(ctx context.Context, userID int) ([]models.ResultSummary, error) {
	panic("unimplemented")
}

// GetResultByQuizNo implements interfaces.ClientInterface.
func (w *appClient) GetResultByQuizNo(ctx context.Context, userID int, QuizNo int) (*models.Result, error) {
	panic("unimplemented")
}

// CreateQuiz implements ClientInterface.
func (w *appClient) CreateQuiz(ctx context.Context, userID int, CreateQuizRequest *models.CreateQuizReq) (*models.CreateQuizResponse, error) {
	panic("unimplemented")
}

// SubmitAnswerAndRoute implements ClientInterface.
func (w *appClient) SubmitAnswerAndRoute(ctx context.Context, userID int, CreateQuizRequest *models.PostAnswerQuestionRequest) (*models.AnswerRouteRes, error) {
	panic("unimplemented")
}

// GetNextOrResume implements ClientInterface.
func (w *appClient) GetNextOrResume(ctx context.Context, userID int, req *models.GetQuizRequest) (*models.GetQuizResponse, error) {
	panic("unimplemented")
}

// BulkRegister implements interfaces.ClientInterface.
func (c *appClient) BulkRegister(ctx context.Context, userID int, words []string) (*models.BulkRegisterResponse, error) {
	panic("unimplemented")
}

// BulkTokenize implements interfaces.ClientInterface.
func (c *appClient) BulkTokenize(ctx context.Context, userID int, text string) ([]string, []string, []string, error) {
	panic("unimplemented")
}

// RootConfig implements interfaces.ClientInterface.
func (c *appClient) RootConfig() *ent.RootConfigClient {
	return c.entClient.RootConfig
}

// UserConfig implements interfaces.ClientInterface.
func (c *appClient) UserConfig() *ent.UserConfigClient {
	return c.entClient.UserConfig
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

// RegisteredWordCount implements interfaces.ClientInterface.
func (c *appClient) RegisteredWordCount(ctx context.Context, RegisteredWordCountRequest *models.RegisteredWordCountRequest) (*models.RegisteredWordCountResponse, error) {
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

// Quiz は QuizClient を返します。
func (c *appClient) Quiz() *ent.QuizClient {
	return c.entClient.Quiz
}

// QuizQuestion は QuizQuestionClient を返します。
func (c *appClient) QuizQuestion() *ent.QuizQuestionClient {
	return c.entClient.QuizQuestion
}

// ExternalAuth は ExternalAuth を返します。
func (c *appClient) ExternalAuth() *ent.ExternalAuthClient {
	return c.entClient.ExternalAuth
}
