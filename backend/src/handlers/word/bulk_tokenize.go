package word

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"word_app/backend/config"
	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) BulkTokenizeHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()
		// 1) 50KB 超を入口で即 413
		limitCfg := config.NewLimitsConfig()
		lr := io.LimitedReader{R: c.Request.Body, N: int64(limitCfg.BulkMaxBytes) + 1}
		body, err := io.ReadAll(&lr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		if len(body) > limitCfg.BulkMaxBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "text too large (limit 50KB)"})
			return
		}

		var req models.BulkTokenizeRequest
		if err := json.Unmarshal(body, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		textBytes := len([]byte(req.Text))
		if textBytes > limitCfg.BulkMaxBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "text too large (limit 50KB)"})
			return
		}

		// 2) 日次回数（5回）を原子的に消費 → 上限で 429
		if _, err := h.userDailyUsageRepo.IncBulkOr429(c, userID, h.clock.Now()); err != nil {
			httperr.Write(c, err) // 429 を返す (ucerr.TooManyRequests)
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
