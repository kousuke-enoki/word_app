package word

import (
	"context"
	"net/http"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func (h *Handler) BulkRegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		// userID は認証ミドルウェアでセットされている前提
		userID, ok := c.Get("userID")
		if !ok {
			logrus.Errorf("userID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
			return
		}
		// userIDの型チェック
		userIDInt, ok := userID.(int)
		if !ok {
			logrus.Errorf("invalid userID type")
			c.JSON(http.StatusBadRequest, gin.H{"errors": "invalid userID type"})
			return
		}

		var req models.BulkRegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			if ve, ok := err.(validator.ValidationErrors); ok {
				errs := make([]gin.H, 0, len(ve))
				for _, fe := range ve {
					errs = append(errs, gin.H{
						"field":   fe.Field(),
						"message": fe.Tag(),
					})
				}
				c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := h.wordService.BulkRegister(ctx, userIDInt, req.Words)
		if err != nil {
			logrus.Errorf("Failed to bulk register: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
