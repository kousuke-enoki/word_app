package interfaces

import (
	"context"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type ExamHandler interface {
	CreateExamHandler() gin.HandlerFunc
}

type ExamService interface {
	CreateExam(ctx context.Context, CreateExamRequest *models.CreateExamRequest) (*models.CreateExamResponse, error)
}
