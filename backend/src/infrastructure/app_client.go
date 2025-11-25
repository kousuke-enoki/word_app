// infrastructure/app_client.go
package infrastructure

import (
	"context"
	"database/sql"

	"word_app/backend/ent"
	"word_app/backend/src/domain"
	"word_app/backend/src/interfaces"
	"word_app/backend/src/models"
	settingUc "word_app/backend/src/usecase/setting"
	user_usecase "word_app/backend/src/usecase/user"
)

// entの型からはmockeryでモックを作れないので、
// こちらをserviceinterfaces.EntClientInterfaceでラップして
// mockeryでモック作成するためのもの
// クリーンアーキテクチャ移行により削除予定

type appClient struct {
	entClient *ent.Client
}

// DB implements interfaces.ClientInterface.
func (c *appClient) DB() *sql.DB {
	panic("unimplemented")
}

// GetDetailByID implements interfaces.ClientInterface.
func (c *appClient) GetDetailByID(ctx context.Context, viewerID int, targetID int) (*models.UserDetail, error) {
	panic("unimplemented")
}

// GetMyDetail implements interfaces.ClientInterface.
func (c *appClient) GetMyDetail(ctx context.Context, viewerID int) (*models.UserDetail, error) {
	panic("unimplemented")
}

// SignUp implements interfaces.ClientInterface.
func (c *appClient) SignUp(ctx context.Context, in user_usecase.SignUpInput) (*user_usecase.SignUpOutput, error) {
	panic("unimplemented")
}

// Update implements interfaces.ClientInterface.
func (c *appClient) FindByEmail(ctx context.Context, email string) (*user_usecase.FindByEmailOutput, error) {
	panic("unimplemented")
}

// Update implements interfaces.ClientInterface.
func (c *appClient) UpdateUser(ctx context.Context, in user_usecase.UpdateUserInput) (*models.UserDetail, error) {
	panic("unimplemented")
}

// GetUsers implements interfaces.ClientInterface.
func (c *appClient) ListUsers(ctx context.Context, in user_usecase.ListUsersInput) (*user_usecase.UserListResponse, error) {
	panic("unimplemented")
}

// Delete implements interfaces.ClientInterface.
func (c *appClient) Delete(ctx context.Context, in user_usecase.DeleteUserInput) error {
	panic("unimplemented")
}

// NewAppClient 初期化関数
func NewAppClient(entClient *ent.Client) interfaces.ClientInterface {
	return &appClient{
		entClient: entClient,
	}
}

// GetRuntimeConfig implements interfaces.ClientInterface.
func (c *appClient) GetRuntimeConfig(_ context.Context) (*settingUc.RuntimeConfigDTO, error) {
	panic("unimplemented")
}

// GetRoot implements interfaces.ClientInterface.
func (c *appClient) GetRoot(_ context.Context, in settingUc.InputGetRootConfig) (*settingUc.OutputGetRootConfig, error) {
	panic("unimplemented")
}

// GetUser implements interfaces.ClientInterface.
func (c *appClient) GetUser(_ context.Context, in settingUc.InputGetUserConfig) (*settingUc.OutputGetUserConfig, error) {
	panic("unimplemented")
}

// UpdateRoot implements interfaces.ClientInterface.
func (c *appClient) UpdateRoot(_ context.Context, in settingUc.InputUpdateRootConfig) (*domain.RootConfig, error) {
	panic("unimplemented")
}

// GetRootConfigExecute implements interfaces.ClientInterface.
func (c *appClient) GetRootConfigExecute(_ context.Context, in settingUc.InputGetRootConfig) (*settingUc.OutputGetRootConfig, error) {
	panic("unimplemented")
}

// GetUserConfigExecute implements interfaces.ClientInterface.
func (c *appClient) GetUserConfigExecute(_ context.Context, in settingUc.InputGetUserConfig) (*settingUc.OutputGetUserConfig, error) {
	panic("unimplemented")
}

// UpdateRootConfigExecute implements interfaces.ClientInterface.
func (c *appClient) UpdateRootConfigExecute(_ context.Context, in settingUc.InputUpdateRootConfig) (*domain.RootConfig, error) {
	panic("unimplemented")
}

// UpdateUserConfigExecute implements interfaces.ClientInterface.
func (c *appClient) UpdateUserConfigExecute(_ context.Context, in settingUc.InputUpdateUserConfig) (*domain.UserConfig, error) {
	panic("unimplemented")
}

// GetSummaries implements interfaces.ClientInterface.
func (c *appClient) GetSummaries(_ context.Context, userID int) ([]models.ResultSummary, error) {
	panic("unimplemented")
}

// GetByQuizNo implements interfaces.ClientInterface.
func (c *appClient) GetByQuizNo(_ context.Context, userID int, QuizNo int) (*models.Result, error) {
	panic("unimplemented")
}

// CreateQuiz implements ClientInterface.
func (c *appClient) CreateQuiz(_ context.Context, userID int, CreateQuizRequest *models.CreateQuizReq) (*models.CreateQuizResponse, error) {
	panic("unimplemented")
}

// SubmitAnswerAndRoute implements ClientInterface.
func (c *appClient) SubmitAnswerAndRoute(_ context.Context, userID int, CreateQuizRequest *models.PostAnswerQuestionRequest) (*models.AnswerRouteRes, error) {
	panic("unimplemented")
}

// GetNextOrResume implements ClientInterface.
func (c *appClient) GetNextOrResume(_ context.Context, userID int, req *models.GetQuizRequest) (*models.GetQuizResponse, error) {
	panic("unimplemented")
}

// BulkRegister implements interfaces.ClientInterface.
func (c *appClient) BulkRegister(_ context.Context, userID int, words []string) (*models.BulkRegisterResponse, error) {
	panic("unimplemented")
}

// BulkTokenize implements interfaces.ClientInterface.
func (c *appClient) BulkTokenize(_ context.Context, userID int, text string) ([]string, []string, []string, error) {
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

// DeleteWord implements interfaces.ClientInterface.
func (c *appClient) DeleteWord(_ context.Context, DeleteWordRequest *models.DeleteWordRequest) (*models.DeleteWordResponse, error) {
	panic("unimplemented")
}

// RegisteredWordCount implements interfaces.ClientInterface.
func (c *appClient) RegisteredWordCount(_ context.Context, RegisteredWordCountRequest *models.RegisteredWordCountRequest) (*models.RegisteredWordCountResponse, error) {
	panic("unimplemented")
}

// GetRegisteredWords implements interfaces.ClientInterface.
func (c *appClient) GetRegisteredWords(_ context.Context, WordListRequest *models.WordListRequest) (*models.WordListResponse, error) {
	panic("unimplemented")
}

// GetWordDetails implements interfaces.ClientInterface.
func (c *appClient) GetWordDetails(_ context.Context, WordShowRequest *models.WordShowRequest) (*models.WordShowResponse, error) {
	panic("unimplemented")
}

// GetWords implements interfaces.ClientInterface.
func (c *appClient) GetWords(_ context.Context, WordListRequest *models.WordListRequest) (*models.WordListResponse, error) {
	panic("unimplemented")
}

// RegisterWords implements interfaces.ClientInterface.
func (c *appClient) RegisterWords(_ context.Context, RegisterWordRequest *models.RegisterWordRequest) (*models.RegisterWordResponse, error) {
	panic("unimplemented")
}

// SaveMemo implements interfaces.ClientInterface.
func (c *appClient) SaveMemo(_ context.Context, SaveMemoRequest *models.SaveMemoRequest) (*models.SaveMemoResponse, error) {
	panic("unimplemented")
}

// UpdateWord implements interfaces.ClientInterface.
func (c *appClient) UpdateWord(_ context.Context, UpdateWordRequest *models.UpdateWordRequest) (*models.UpdateWordResponse, error) {
	panic("unimplemented")
}

// UserClient の実装
func (c *appClient) Create(ctx context.Context, email, name, password string) (*ent.User, error) {
	return c.entClient.User.Create().
		SetEmail(email).
		SetName(name).
		SetPassword(password).
		Save(ctx)
}

// func (c *appClient) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
// 	return c.entClient.User.Query().Where(user.Email(email)).Only(ctx)
// }

// func (c *appClient) FindByID(ctx context.Context, id int) (*ent.User, error) {
// 	return c.entClient.User.Query().Where(user.ID(id)).Only(ctx)
// }

// WordService の実装
func (c *appClient) CreateWord(_ context.Context, req *models.CreateWordRequest) (*models.CreateWordResponse, error) {
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

// UserDailyUsage implements interfaces.ClientInterface.
func (c *appClient) UserDailyUsage() *ent.UserDailyUsageClient {
	return c.entClient.UserDailyUsage
}
