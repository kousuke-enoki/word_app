// handlers/user/detail.go
package user

import (
	"net/http"
	"strconv"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) MeHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, viewerID int) {
		ctx := c.Request.Context()

		dto, err := h.userUsecase.GetMyDetail(ctx, viewerID)
		if err != nil {
			httperr.Write(c, err)
			return
		}
		c.JSON(http.StatusOK, dto)
	})
}

func (h *UserHandler) ShowHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, viewerID int) {
		ctx := c.Request.Context()

		targetID, parseErr := strconv.Atoi(c.Param("id"))
		if parseErr != nil || targetID <= 0 {
			httperr.Write(c, apperror.New(apperror.Validation, "invalid id", parseErr))
			return
		}

		dto, err := h.userUsecase.GetDetailByID(ctx, viewerID, targetID)
		if err != nil {
			httperr.Write(c, err)
			return
		}
		c.JSON(http.StatusOK, dto)
	})
}
