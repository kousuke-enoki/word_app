package interfaces

import (
	"context"
	"word_app/backend/ent"
	"word_app/backend/src/domain"
	"word_app/backend/src/models"
	settingUc "word_app/backend/src/usecase/setting"
	"word_app/backend/src/utils/contextutil"
)

type TestClientWrapper struct {
	entClient *ent.Client
}

// GetRootConfigExecute implements ClientInterface.
func (w *TestClientWrapper) GetRootConfigExecute(ctx context.Context, in settingUc.GetRootConfigInput) (*settingUc.GetRootConfigOutput, error) {
	panic("unimplemented")
}

// GetUserConfigExecute implements ClientInterface.
func (w *TestClientWrapper) GetUserConfigExecute(ctx context.Context, in settingUc.GetUserConfigInput) (*settingUc.GetUserConfigOutput, error) {
	panic("unimplemented")
}

// UpdateRootConfigExecute implements ClientInterface.
func (w *TestClientWrapper) UpdateRootConfigExecute(ctx context.Context, in settingUc.UpdateRootConfigInput) (*domain.RootConfig, error) {
	panic("unimplemented")
}

// UpdateUserConfigExecute implements ClientInterface.
func (w *TestClientWrapper) UpdateUserConfigExecute(ctx context.Context, in settingUc.UpdateUserConfigInput) (*domain.UserConfig, error) {
	panic("unimplemented")
}

// GetAuthConfig implements interfaces.ClientInterface.
func (w *TestClientWrapper) GetAuthConfig(ctx context.Context) (*settingUc.AuthConfigDTO, error) {
	panic("unimplemented")
}

// GetAuthConfig implements interfaces.ClientInterface.
func (w *TestClientWrapper) GetAuthConfigs(ctx context.Context) (*models.AuthSettingResponse, error) {
	panic("unimplemented")
}

// GetResultSummaries implements ClientInterface.
func (w *TestClientWrapper) Validate(ctx context.Context, tokenStr string) (contextutil.UserRoles, error) {
	panic("unimplemented")
}

// GetResultSummaries implements ClientInterface.
func (w *TestClientWrapper) GetResultSummaries(ctx context.Context, userID int) ([]models.ResultSummary, error) {
	panic("unimplemented")
}

// GetResultByQuizNo implements ClientInterface.
func (w *TestClientWrapper) GetResultByQuizNo(ctx context.Context, userID int, QuizNo int) (*models.Result, error) {
	panic("unimplemented")
}

// CreateQuiz implements ClientInterface.
func (w *TestClientWrapper) CreateQuiz(ctx context.Context, userID int, CreateQuizRequest *models.CreateQuizReq) (*models.CreateQuizResponse, error) {
	panic("unimplemented")
}

// SubmitAnswerAndRoute implements ClientInterface.
func (w *TestClientWrapper) SubmitAnswerAndRoute(ctx context.Context, userID int, CreateQuizRequest *models.PostAnswerQuestionRequest) (*models.AnswerRouteRes, error) {
	panic("unimplemented")
}

// finishQuizTx implements ClientInterface.
func (w *TestClientWrapper) finishQuizTx(ctx context.Context, tx *ent.Tx, q *ent.Quiz) (*models.Result, error) {
	panic("unimplemented")
}

// GetNextOrResume implements ClientInterface.
func (w *TestClientWrapper) GetNextOrResume(ctx context.Context, userID int, req *models.GetQuizRequest) (*models.GetQuizResponse, error) {
	panic("unimplemented")
}

// BulkRegister implements ClientInterface.
func (w *TestClientWrapper) BulkRegister(ctx context.Context, userID int, words []string) (*models.BulkRegisterResponse, error) {
	panic("unimplemented")
}

// BulkTokenize implements ClientInterface.
func (w *TestClientWrapper) BulkTokenize(ctx context.Context, userID int, text string) ([]string, []string, []string, error) {
	panic("unimplemented")
}

// GetRootConfig implements ClientInterface.
func (w *TestClientWrapper) GetRootConfig(ctx context.Context, userID int) (*ent.RootConfig, error) {
	panic("unimplemented")
}

// GetUserConfig implements ClientInterface.
func (w *TestClientWrapper) GetUserConfig(ctx context.Context, userID int) (*ent.UserConfig, error) {
	panic("unimplemented")
}

// UpdateRootConfig implements ClientInterface.
func (w *TestClientWrapper) UpdateRootConfig(ctx context.Context, userID int, editingPermissions string, isTestUserMode bool, IsEmailAuthCheck bool, isLineAuth bool) (*ent.RootConfig, error) {
	panic("unimplemented")
}

// UpdateUserConfig implements ClientInterface.
func (w *TestClientWrapper) UpdateUserConfig(ctx context.Context, userID int, isLightMode bool) (*ent.UserConfig, error) {
	panic("unimplemented")
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

// RegisteredWordCount implements ClientInterface
func (w *TestClientWrapper) RegisteredWordCount(ctx context.Context, RegisteredWordCountRequest *models.RegisteredWordCountRequest) (*models.RegisteredWordCountResponse, error) {
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

// User は UserClient を返します。
func (w *TestClientWrapper) UserConfig() *ent.UserConfigClient {
	panic("unimplemented")
}

// User は UserClient を返します。
func (w *TestClientWrapper) RootConfig() *ent.RootConfigClient {
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

// Quiz は QuizClient を返します。
func (w *TestClientWrapper) Quiz() *ent.QuizClient {
	panic("unimplemented")
}

// QuizQuestion は QuizQuestionClient を返します。
func (w *TestClientWrapper) QuizQuestion() *ent.QuizQuestionClient {
	panic("unimplemented")
}

// ExternalAuth は ExternalAuthClient を返します。
func (w *TestClientWrapper) ExternalAuth() *ent.ExternalAuthClient {
	panic("unimplemented")
}
