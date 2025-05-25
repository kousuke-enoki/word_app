package interfaces

import (
	"context"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type QuizHandler interface {
	CreateQuizHandler() gin.HandlerFunc
	PostAnswerAndRouteHandler() gin.HandlerFunc
	GetQuizHandler() gin.HandlerFunc
}

type QuizService interface {
	CreateQuiz(ctx context.Context, userID int, CreateQuizRequest *models.CreateQuizReq) (*models.CreateQuizResponse, error)
	SubmitAnswerAndRoute(ctx context.Context, userID int, in *models.PostAnswerQuestionRequest) (*models.AnswerRouteRes, error)
	GetNextOrResume(ctx context.Context, userID int, req *models.GetQuizRequest) (*models.GetQuizResponse, error)
}
