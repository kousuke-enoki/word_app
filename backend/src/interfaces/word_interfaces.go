package interfaces

import (
	"context"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type WordHandler interface {
	WordNewHandler() gin.HandlerFunc
	DeleteWordHandler() gin.HandlerFunc
	AllWordListHandler() gin.HandlerFunc
	WordShowHandler() gin.HandlerFunc
	RegisterWordHandler() gin.HandlerFunc
	SaveMemoHandler() gin.HandlerFunc
}

type WordService interface {
	CreateWord(ctx context.Context, WordCreateRequest *models.CreateWordRequest) (*models.CreateWordResponse, error)
	GetWordDetails(ctx context.Context, wordID int, userID int) (*models.WordShowResponse, error)
	GetWords(ctx context.Context, userID int, search string, sortBy string, order string, page int, limit int) (*models.AllWordListResponse, error)
	DeleteWord(ctx context.Context, userID int, wordID int) (*models.DeleteWordResponse, error)
	GetRegisteredWords(ctx context.Context, userID int, search string, order string, page int, limit int) (*models.AllWordListResponse, error)
	RegisterWords(ctx context.Context, wordID int, userID int, IsRegistered bool) (*models.RegisterWordResponse, error)
	SaveMemo(ctx context.Context, wordID int, userID int, memo string) (*models.SaveMemoResponse, error)
}
