// handlers/user/detail.go
package user

import (
	"net/http"
	"strconv"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/usecase/apperror"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) MeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		viewerID, err := contextutil.MustUserID(c)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		dto, err := h.userUsecase.GetMyDetail(c.Request.Context(), viewerID)
		if err != nil {
			httperr.Write(c, err)
			return
		}
		c.JSON(http.StatusOK, dto)
	}
}

func (h *UserHandler) ShowHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		viewerID, err := contextutil.MustUserID(c)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		targetID, parseErr := strconv.Atoi(c.Param("id"))
		if parseErr != nil || targetID <= 0 {
			httperr.Write(c, apperror.New(apperror.Validation, "invalid id", parseErr))
			return
		}

		dto, err := h.userUsecase.GetDetailByID(c.Request.Context(), viewerID, targetID)
		if err != nil {
			httperr.Write(c, err)
			return
		}
		c.JSON(http.StatusOK, dto)
	}
}
