package result

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GET /results
func (h *ResultHandler) GetResultsIndexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		uidRaw, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, ok := uidRaw.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid userID"})
			return
		}

		list, err := h.resultService.GetResultSummaries(ctx, userID)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch results"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}
