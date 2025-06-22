// jwt_check_middleware_test.go
package jwt

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/src/middleware/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestJwtCheckMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 共通: JSON 比較 helper
	// trim := func(b []byte) string { return string(bytes.TrimSpace(b)) }

	t.Run("roles present → 200", func(t *testing.T) {
		r := gin.New()

		// ① ロールを仕込むダミーミドルウェア
		r.Use(func(c *gin.Context) {
			c.Set("userID", 42)
			c.Set("isAdmin", true)
			c.Set("isRoot", false)
		})
		// ② テスト対象
		r.GET("/mypage", new(jwt.JwtMiddleware).JwtCheckMiddleware())

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/mypage", nil)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		// want := `{"user":{"id":42,"isAdmin":true,"isRoot":false}}`
		require.JSONEq(t,
			`{"user":{"id":42,"name":"","isAdmin":true,"isRoot":false}}`, // ← name を追加
			string(bytes.TrimSpace(w.Body.Bytes())),
		)
	})

	t.Run("roles missing → 401", func(t *testing.T) {
		r := gin.New()
		// ロールを仕込まない
		r.GET("/mypage", new(jwt.JwtMiddleware).JwtCheckMiddleware())

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/mypage", nil)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		require.Contains(t, w.Body.String(), "userID not found in context")
	})
}
