package interfaces

import (
	"context"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type ResultHandler interface {
	GetResultsIndexHandler() gin.HandlerFunc
	GetResultHandler() gin.HandlerFunc
}

type ResultService interface {
	GetResultSummaries(ctx context.Context, userID int) ([]models.ResultSummary, error)
	GetResultByQuizNo(ctx context.Context, userID int, QuizNo int) (*models.Result, error)
}
