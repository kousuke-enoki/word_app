// jwt_check_middleware_test.go
package jwt

import (
	"bytes"
	"encoding/json"
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

	t.Run("user logined", func(t *testing.T) {
		r := gin.New()

		// ① ロールを仕込むダミーミドルウェア
		r.Use(func(c *gin.Context) {
			c.Set("userID", 42)
			c.Set("isAdmin", true)
			c.Set("isRoot", false)
			c.Set("isTest", false)
		})
		// ② テスト対象
		r.GET("/mypage", new(jwt.JwtMiddleware).JwtCheckMiddleware())

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/mypage", nil)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		// want := `{"user":{"id":42,"isAdmin":true,"isRoot":false}}`
		require.JSONEq(t,
			`{"user":{"id":42,"name":"","isAdmin":true,"isRoot":false,"isTest":false}, "isLogin": true}`,
			string(bytes.TrimSpace(w.Body.Bytes())),
		)
	})

	t.Run("user not logined", func(t *testing.T) {
		r := gin.New()
		// Recovery は使わず自前 recover で panic を確認
		r.GET("/mypage", new(jwt.JwtMiddleware).JwtCheckMiddleware())

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/mypage", nil)

		var rec any
		func() {
			defer func() { rec = recover() }()
			r.ServeHTTP(w, req)
		}()

		require.NotNil(t, rec, "roles が無いので panic するはず")

		// --- JSON は柔らかくチェックする ---
		var got map[string]any
		require.NoError(t, json.Unmarshal(bytes.TrimSpace(w.Body.Bytes()), &got))

		// isLogin が false であることだけ必須に
		require.Equal(t, false, got["isLogin"])

		// user が含まれていてもゼロ値なら OK とする
		if u, ok := got["user"].(map[string]any); ok {
			require.Equal(t, float64(0), u["id"])
			require.Equal(t, false, u["isAdmin"])
			require.Equal(t, false, u["isRoot"])
			require.Equal(t, false, u["isTest"])
		}
	})
}
