package jwt

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 後続ハンドラでユーザーロールを使えるようにcontextにセットするミドルウェア
// ユーザーID、isAdmin、isRootをcontextにセットする
func (m *JwtMiddleware) AuthenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}
		p, err := m.JwtUsecase.Authenticate(c.Request.Context(), raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		// Principal を context へ（キーは型安全に）
		SetPrincipal(c, p)
		// アクセスログ用にuser_idを設定
		c.Set("user_id", p.UserID)
		// 標準 context にも入れて、Request に戻す（以降は c.Request.Context() で拾える）
		ctx := WithPrincipal(c.Request.Context(), p)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
