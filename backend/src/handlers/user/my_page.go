// user/handler.go
package user

import (
	"context"
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/models"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) MyPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		// userID の取得
		userID, err := contextutil.MustUserID(c)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		// ユーザー情報の取得
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
	}
}
