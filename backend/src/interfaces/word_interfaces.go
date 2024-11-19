package interfaces

import (
	"context"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type WordHandler interface {
	AllWordListHandler() gin.HandlerFunc
	WordShowHandler() gin.HandlerFunc
}

type WordService interface {
	GetWordDetails(ctx context.Context, wordID int) (*models.WordResponse, error)
	GetWords(ctx context.Context, search string, sortBy string, order string, page int, limit int) ([]models.Word, int, int, error)
}
