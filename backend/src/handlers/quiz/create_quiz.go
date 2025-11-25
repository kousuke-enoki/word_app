package quiz

import (
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()
		req, err := h.parseCreateQuizRequest(c)
		if err != nil {
			logrus.Errorf("Failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := h.quizService.CreateQuiz(ctx, userID, req)
		if err != nil {
			httperr.Write(c, err) // apperrorをそのまま返す（429も含む）
			return
		}

		c.JSON(http.StatusOK, response)
	})
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
