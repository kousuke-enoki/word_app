package jwt

import (
	"net/http"
	"word_app/backend/src/models"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
)

// JwtAuth : JWT検証 & ユーザー情報(ロール)取得
// GinのミドルウェアとしてJWTを検証し、ユーザーロールをフロントに返す
func (m *JwtMiddleware) JwtCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		roles, err := contextutil.GetUserRoles(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.MyPageResponse{
			User: models.User{
				ID:      roles.UserID,
				IsAdmin: roles.IsAdmin,
				IsRoot:  roles.IsRoot,
			},
		})
	}
}
