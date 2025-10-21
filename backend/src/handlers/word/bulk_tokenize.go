package word

import (
	"net/http"
	"strings"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) BulkTokenizeHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()

		var req models.BulkTokenizeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"errors": err})
			return
		}

		cands, regs, notExist, err := h.wordService.BulkTokenize(ctx, userID, req.Text)
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
	})
}
