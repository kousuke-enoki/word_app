package result

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *ResultHandler) GetResultHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*--- userID は認証ミドルウェアで埋め込んである想定 ---*/
		rawID, ok := c.Get("userID")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID := rawID.(int)

		/*--- パスパラメータ ---*/
		noStr := c.Param("quizNo")
		quizNo, err := strconv.Atoi(noStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quizNo"})
			return
		}

		/*--- サービス呼び出し ---*/
		res, err := h.resultService.GetResultByQuizNo(c.Request.Context(), userID, quizNo)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "result not found"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
