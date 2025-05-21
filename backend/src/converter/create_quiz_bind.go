package converter

import (
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

func BindCreateQuiz(c *gin.Context) (*models.CreateQuizDTO, error) {
	var dto models.CreateQuizDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		return nil, err
	}
	return &dto, nil
}
