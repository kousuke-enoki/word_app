package bulk

import (
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
)

func (h *BulkHandler) RegisterHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		var req models.BulkRegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			httperr.Write(c, apperror.BadRequestf("invalid json", err))
			return
		}
		res, err := h.registerUsecase.Register(c, userID, req.Words)
		if err != nil {
			httperr.Write(c, err) // apperrorをそのまま返す
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
