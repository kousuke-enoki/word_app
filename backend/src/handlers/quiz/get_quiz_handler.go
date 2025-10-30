package quiz

import (
	"net/http"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()
		// --- query パラメータをパース ---
		var req models.GetQuizRequest

		// --- サービス呼び出し ---
		q, err := h.quizService.GetNextOrResume(ctx, userID, &req)
		if err != nil && q == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, q)
	})
}
