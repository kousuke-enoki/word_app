package interfaces

import (
	"context"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type QuizHandler interface {
	CreateQuizHandler() gin.HandlerFunc
}

type QuizService interface {
	CreateQuiz(ctx context.Context, CreateQuizRequest *models.CreateQuizRequest) (*models.CreateQuizResponse, error)
}
