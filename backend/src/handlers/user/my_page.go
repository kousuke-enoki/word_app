// user/handler.go
package user

import (
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) MyPageHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		// ← ここは userID が必ずある前提の世界
		ctx := c.Request.Context()
		signInUser, err := h.userUsecase.GetMyDetail(ctx, userID)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		c.JSON(http.StatusOK, models.MyPageResponse{
			User: models.User{
				ID:      signInUser.ID,
				Name:    signInUser.Name,
				IsAdmin: signInUser.IsAdmin,
				IsRoot:  signInUser.IsRoot,
				IsTest:  signInUser.IsTest,
			},
			IsLogin: true,
		})
	})
}
