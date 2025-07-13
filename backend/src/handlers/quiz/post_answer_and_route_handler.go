package quiz

import (
	"context"
	"net/http"

	"word_app/backend/src/models"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *QuizHandler) PostAnswerAndRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		// リクエストを解析
		req, err := h.parsePostAnswerQuizRequest(c)
		if err != nil {
			logrus.Errorf("Failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logrus.Debug(req)

		// ユーザーIDをコンテキストから取得
		userID, err := contextutil.MustUserID(c)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
	}
}

// リクエスト構造体を解析
func (h *QuizHandler) parsePostAnswerQuizRequest(c *gin.Context) (*models.PostAnswerQuestionRequest, error) {
	var req models.PostAnswerQuestionRequest

	// JSONリクエストをバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		return nil, err
	}

	return &req, nil
}
