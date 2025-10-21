package result

import (
	"net/http"

	"word_app/backend/src/middleware/jwt"

	"github.com/gin-gonic/gin"
)

// GET /results
func (h *Handler) GetIndexHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()

		list, err := h.resultService.GetSummaries(ctx, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch results"})
			return
		}
		c.JSON(http.StatusOK, list)
	})
}
