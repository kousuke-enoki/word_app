package quiz

import (
	"net/http"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) PostAnswerAndRouteHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()
		// リクエストを解析
		req, err := h.parsePostAnswerQuizRequest(c)
		if err != nil {
			logrus.Errorf("Failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// サービス層にリクエストを渡して処理
		response, err := h.quizService.SubmitAnswerAndRoute(ctx, userID, req)
		if err != nil {
			logrus.Errorf("Failed to submit answer: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit answer"})
			return
		}

		c.JSON(http.StatusOK, response)
	})
}

// リクエスト構造体を解析
func (h *Handler) parsePostAnswerQuizRequest(c *gin.Context) (*models.PostAnswerQuestionRequest, error) {
	var req models.PostAnswerQuestionRequest

	// JSONリクエストをバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		return nil, err
	}

	return &req, nil
}
