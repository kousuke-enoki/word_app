package quiz

import (
	"context"
	"net/http"

	"word_app/backend/src/models"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		req, err := h.parseCreateQuizRequest(c)
		if err != nil {
			logrus.Errorf("Failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logrus.Debug(req)

		userID, err := contextutil.MustUserID(c)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		response, err := h.quizService.CreateQuiz(ctx, userID, req)
		if err != nil {
			logrus.Errorf("Failed to create quiz: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz"})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// リクエスト構造体を解析
func (h *Handler) parseCreateQuizRequest(c *gin.Context) (*models.CreateQuizReq, error) {
	var req models.CreateQuizReq

	// JSONリクエストをバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		return nil, err
	}

	return &req, nil
}
