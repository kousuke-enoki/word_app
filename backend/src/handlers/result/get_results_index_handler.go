package result

import (
	"context"
	"net/http"

	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GET /results
func (h *Handler) GetIndexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		userID, err := contextutil.MustUserID(c)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		list, err := h.resultService.GetSummaries(ctx, userID)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch results"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}
