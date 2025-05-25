package result

import (
	"net/http"
	"strconv"
	"word_app/backend/src/handlers/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *ResultHandler) GetResultHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.MustUserID(c)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		noStr := c.Param("quizNo")
		quizNo, err := strconv.Atoi(noStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quizNo"})
			return
		}
		logrus.Debug(quizNo)

		res, err := h.resultService.GetResultByQuizNo(c.Request.Context(), userID, quizNo)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "result not found"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
