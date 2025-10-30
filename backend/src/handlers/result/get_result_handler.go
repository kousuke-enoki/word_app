package result

import (
	"net/http"
	"strconv"

	"word_app/backend/src/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()

		noStr := c.Param("quizNo")
		quizNo, err := strconv.Atoi(noStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quizNo"})
			return
		}

		res, err := h.resultService.GetByQuizNo(ctx, userID, quizNo)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "result not found"})
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
