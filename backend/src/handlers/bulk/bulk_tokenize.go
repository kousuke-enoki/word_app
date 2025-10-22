package word

import (
	"encoding/json"
	"io"
	"net/http"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

func (h *BulkHandler) Tokenize() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		// 入口で50KB超を即413
		lr := io.LimitedReader{R: c.Request.Body, N: int64(h.limits.BulkMaxBytes) + 1}
		body, err := io.ReadAll(&lr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		if len(body) > h.limits.BulkMaxBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "text too large (limit 50KB)"})
			return
		}
		var req models.BulkTokenizeRequest
		if err := json.Unmarshal(body, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		cands, regs, notExist, err := h.tokenizeUsecase.Execute(c, userID, req.Text)
		if err != nil {
			// ucerr.TooManyRequests(429) 等をそのまま返す（あなたの httperr があれば使う）
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, models.BulkTokenizeResponse{Candidates: cands, Registered: regs, NotExistWord: notExist})
	})
}
