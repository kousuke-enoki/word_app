package word

import (
	"context"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateHandler() gin.HandlerFunc
	UpdateHandler() gin.HandlerFunc
	DeleteHandler() gin.HandlerFunc
	ListHandler() gin.HandlerFunc
	ShowHandler() gin.HandlerFunc
	RegisterHandler() gin.HandlerFunc
	SaveMemoHandler() gin.HandlerFunc
	BulkTokenizeHandler() gin.HandlerFunc
	BulkRegisterHandler() gin.HandlerFunc
}

type Service interface {
	CreateWord(ctx context.Context, CreateWordRequest *models.CreateWordRequest) (*models.CreateWordResponse, error)
	UpdateWord(ctx context.Context, UpdateWordRequest *models.UpdateWordRequest) (*models.UpdateWordResponse, error)
	GetWordDetails(ctx context.Context, WordShowRequest *models.WordShowRequest) (*models.WordShowResponse, error)
	GetWords(ctx context.Context, WordListRequest *models.WordListRequest) (*models.WordListResponse, error)
	DeleteWord(ctx context.Context, DeleteWordRequest *models.DeleteWordRequest) (*models.DeleteWordResponse, error)
	GetRegisteredWords(ctx context.Context, WordListRequest *models.WordListRequest) (*models.WordListResponse, error)
	RegisterWords(ctx context.Context, RegisterWordRequest *models.RegisterWordRequest) (*models.RegisterWordResponse, error)
	RegisteredWordCount(ctx context.Context, RegisteredWordCountRequest *models.RegisteredWordCountRequest) (*models.RegisteredWordCountResponse, error)
	SaveMemo(ctx context.Context, SaveMemoRequest *models.SaveMemoRequest) (*models.SaveMemoResponse, error)
	BulkTokenize(ctx context.Context, userID int, text string) ([]string, []string, []string, error)
	BulkRegister(ctx context.Context, userID int, words []string) (*models.BulkRegisterResponse, error)
}
type Validator interface {
	ValidateCreateWordRequest(CreateWordRequest *models.CreateWordRequest) []*models.FieldError
	ValidateSaveMemo(SaveMemoRequest *models.SaveMemoRequest) []*models.FieldError
	ValidateUpdateWordRequest(UpdateWordRequest *models.UpdateWordRequest) []*models.FieldError
}
