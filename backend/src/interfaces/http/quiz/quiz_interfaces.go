package quiz

import (
	"context"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

// Handler bundles the Gin handlers that expose quiz-related HTTP endpoints.
// Each method returns a `gin.HandlerFunc` ready to be registered on a router.
type Handler interface {
	CreateHandler() gin.HandlerFunc             // POST /quizzes
	PostAnswerAndRouteHandler() gin.HandlerFunc // POST /quizzes/:id/answers
	GetHandler() gin.HandlerFunc                // GET  /quizzes/next-or-resume
}

// Service defines the business-logic layer used by the quiz HTTP handlers.
// It hides persistence details behind an interface, making unit-testing easier.
type Service interface {
	CreateQuiz(ctx context.Context, userID int, CreateQuizRequest *models.CreateQuizReq) (*models.CreateQuizResponse, error)
	SubmitAnswerAndRoute(ctx context.Context, userID int, in *models.PostAnswerQuestionRequest) (*models.AnswerRouteRes, error)
	GetNextOrResume(ctx context.Context, userID int, req *models.GetQuizRequest) (*models.GetQuizResponse, error)
}
