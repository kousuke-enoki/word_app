package result

import (
	"context"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetIndexHandler() gin.HandlerFunc
	GetHandler() gin.HandlerFunc
}

type Service interface {
	GetSummaries(ctx context.Context, userID int) ([]models.ResultSummary, error)
	GetByQuizNo(ctx context.Context, userID int, QuizNo int) (*models.Result, error)
}
