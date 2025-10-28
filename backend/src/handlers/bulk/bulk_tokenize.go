package bulk

import (
	"encoding/json"
	"io"
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
)

func (h *BulkHandler) TokenizeHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()
		// 入口で50KB超を即413
		lr := io.LimitedReader{R: c.Request.Body, N: int64(h.limits.BulkMaxBytes) + 1}
		body, err := io.ReadAll(&lr)
		if err != nil {
			httperr.Write(c, apperror.BadRequestf("invalid body", err))
			return
		}
		if len(body) > h.limits.BulkMaxBytes {
			httperr.Write(c, apperror.TooLargeRequestsf("text too large (limit 50KB)", nil))
			return
		}
		var req models.BulkTokenizeRequest
		if err := json.Unmarshal(body, &req); err != nil {
			httperr.Write(c, apperror.BadRequestf("invalid json", err))
			return
		}

		cands, regs, notExist, err := h.tokenizeUsecase.Execute(ctx, userID, req.Text)
		if err != nil {
			httperr.Write(c, err) // apperrorをそのまま返す（429も含む）
			return
		}
		c.JSON(http.StatusOK, models.BulkTokenizeResponse{Candidates: cands, Registered: regs, NotExistWord: notExist})
	})
}
