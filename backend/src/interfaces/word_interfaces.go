package interfaces

import (
	"context"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type WordHandler interface {
	CreateWordHandler() gin.HandlerFunc
	UpdateWordHandler() gin.HandlerFunc
	DeleteWordHandler() gin.HandlerFunc
	WordListHandler() gin.HandlerFunc
	WordShowHandler() gin.HandlerFunc
	RegisterWordHandler() gin.HandlerFunc
	SaveMemoHandler() gin.HandlerFunc
}

type WordService interface {
	CreateWord(ctx context.Context, CreateWordRequest *models.CreateWordRequest) (*models.CreateWordResponse, error)
	UpdateWord(ctx context.Context, UpdateWordRequest *models.UpdateWordRequest) (*models.UpdateWordResponse, error)
	GetWordDetails(ctx context.Context, WordShowRequest *models.WordShowRequest) (*models.WordShowResponse, error)
	GetWords(ctx context.Context, WordListRequest *models.WordListRequest) (*models.WordListResponse, error)
	DeleteWord(ctx context.Context, DeleteWordRequest *models.DeleteWordRequest) (*models.DeleteWordResponse, error)
	GetRegisteredWords(ctx context.Context, WordListRequest *models.WordListRequest) (*models.WordListResponse, error)
	RegisterWords(ctx context.Context, RegisterWordRequest *models.RegisterWordRequest) (*models.RegisterWordResponse, error)
	RegisteredWordCount(ctx context.Context, RegisteredWordCountRequest *models.RegisteredWordCountRequest) (*models.RegisteredWordCountResponse, error)
	SaveMemo(ctx context.Context, SaveMemoRequest *models.SaveMemoRequest) (*models.SaveMemoResponse, error)
}
type WordValidator interface {
	ValidateCreateWordRequest(CreateWordRequest *models.CreateWordRequest) []*models.FieldError
	ValidateSaveMemo(SaveMemoRequest *models.SaveMemoRequest) []*models.FieldError
	ValidateUpdateWordRequest(UpdateWordRequest *models.UpdateWordRequest) []*models.FieldError
}
