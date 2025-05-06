package quiz

import (
	"net/http"
	"strconv"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

func (h *QuizHandler) GetQuizHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// --- 認証済み userID を取得 ---
		raw, ok := c.Get("userID")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, _ := raw.(int)

		// --- query パラメータをパース ---
		var req models.GetQuizRequest

		if v := c.Query("quizId"); v != "" {
			id, _ := strconv.Atoi(v)
			req.QuizID = &id
		}
		if v := c.Query("questionNumber"); v != "" {
			n, _ := strconv.Atoi(v)
			req.BeforeQuestionNumber = &n
		}

		// --- サービス呼び出し ---
		q, err := h.quizService.GetNextOrResume(c.Request.Context(), userID, &req)
		if err != nil {
			// if errors.Is(err, ent.ISNotFound) {
			// 	c.Status(http.StatusNoContent)
			// 	return
			// }
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, q)
	}
}
