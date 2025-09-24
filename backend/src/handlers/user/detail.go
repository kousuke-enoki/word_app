// handlers/user/detail.go
package user

import (
	"net/http"
	"strconv"

	"word_app/backend/src/usecase"

	"github.com/gin-gonic/gin"
)

type DetailHandler struct {
	UC *usecase.UserDetailUsecase
}

func NewDetailHandler(uc *usecase.UserDetailUsecase) *DetailHandler {
	return &DetailHandler{UC: uc}
}

func (h *DetailHandler) MeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		viewerID := v.(int)

		dto, status, err := h.UC.GetMyDetail(c.Request.Context(), viewerID)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, dto)
	}
}

func (h *DetailHandler) ShowHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		viewerIDRaw, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		viewerID, _ := viewerIDRaw.(int)

		idStr := c.Param("id")
		targetID, err := strconv.Atoi(idStr)
		if err != nil || targetID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		dto, err := h.UC.GetDetailByID(c.Request.Context(), viewerID, targetID)
		if err != nil {
			switch err {
			case usecase.ErrUnauthorized:
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			case ErrUserNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, dto)
	}
}
