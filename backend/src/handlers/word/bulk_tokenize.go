package word

import (
	"context"
	"net/http"
	"strings"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) BulkTokenizeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		// userID は認証ミドルウェアでセットされている前提
		userID, ok := c.Get("userID")
		if !ok {
			logrus.Errorf("userID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
			return
		}
		// userIDの型チェック
		userIDInt, ok := userID.(int)
		if !ok {
			logrus.Errorf("invalid userID type")
			c.JSON(http.StatusBadRequest, gin.H{"errors": "invalid userID type"})
			return
		}

		var req models.BulkTokenizeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"errors": err})
			return
		}

		cands, regs, notExist, err := h.wordService.BulkTokenize(ctx, userIDInt, req.Text)
		if err != nil {
			if strings.HasPrefix(err.Error(), "too many tokens") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.BulkTokenizeResponse{
			Candidates:   cands,
			Registered:   regs,
			NotExistWord: notExist,
		})
	}
}
