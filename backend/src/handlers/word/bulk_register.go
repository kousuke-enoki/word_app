package word

import (
	"net/http"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func (h *Handler) BulkRegisterHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()

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

		response, err := h.wordService.BulkRegister(ctx, userID, req.Words)
		if err != nil {
			logrus.Errorf("Failed to bulk register: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	})
}
