// handlers/user/detail.go
package user

import (
	"net/http"
	"strconv"

	"word_app/backend/src/handlers"
	"word_app/backend/src/interfaces/http/user"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
)

type DetailHandler struct {
	UserUC user.Usecase
}

func NewDetailHandler(
	uc user.Usecase,
) *DetailHandler {
	return &DetailHandler{
		UserUC: uc,
	}
}

func (h *DetailHandler) MeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("userID")
		if !ok {
			handlers.WriteError(c, apperror.New(apperror.Unauthorized, "unauthorized", nil))
			return
		}
		viewerID := v.(int)

		dto, err := h.UserUC.GetMyDetail(c.Request.Context(), viewerID)
		if err != nil {
			handlers.WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto)
	}
}

func (h *DetailHandler) ShowHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("userID")
		if !ok {
			handlers.WriteError(c, apperror.New(apperror.Unauthorized, "unauthorized", nil))
			return
		}
		viewerID := v.(int)

		targetID, parseErr := strconv.Atoi(c.Param("id"))
		if parseErr != nil || targetID <= 0 {
			handlers.WriteError(c, apperror.New(apperror.Validation, "invalid id", parseErr))
			return
		}

		dto, err := h.UserUC.GetDetailByID(c.Request.Context(), viewerID, targetID)
		if err != nil {
			handlers.WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto)
	}
}
