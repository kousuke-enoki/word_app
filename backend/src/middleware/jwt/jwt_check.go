package jwt

import (
	"net/http"

	"word_app/backend/src/models"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
)

// JwtAuth : JWT検証 & ユーザー情報(ロール)取得
// GinのミドルウェアとしてJWTを検証し、ユーザーロールをフロントに返す
func (m *Middleware) JwtCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		roles, err := contextutil.GetUserRoles(c)
		if err != nil || roles == nil {
			c.JSON(http.StatusOK, models.MyPageResponse{
				IsLogin: false,
			})
		}

		c.JSON(http.StatusOK, models.MyPageResponse{
			User: models.User{
				ID:      roles.UserID,
				IsAdmin: roles.IsAdmin,
				IsRoot:  roles.IsRoot,
				IsTest:  roles.IsTest,
			},
			IsLogin: true,
		})
	}
}
