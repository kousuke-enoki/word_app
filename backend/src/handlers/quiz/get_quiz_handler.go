package quiz

import (
	"net/http"

	"word_app/backend/src/models"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *QuizHandler) GetQuizHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ユーザーIDをコンテキストから取得
		userID, err := contextutil.MustUserID(c)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// --- query パラメータをパース ---
		var req models.GetQuizRequest
		logrus.Debug(req)

		// --- サービス呼び出し ---
		q, err := h.quizService.GetNextOrResume(c.Request.Context(), userID, &req)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, q)
	}
}
