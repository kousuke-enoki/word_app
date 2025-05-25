package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"word_app/backend/database"
	"word_app/backend/ent/enttest"
	"word_app/backend/src/handlers/middleware"
	"word_app/backend/src/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// 1. テスト用のEntクライアント作成（SQLite in-memory）
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// 2. DBクライアントを差し替え
	database.SetEntClient(client)
	t.Run("TestAuthMiddleware_ValidToken", func(t *testing.T) {
		// 2. DBクライアントを差し替え
		database.SetEntClient(client)

		// 3. テスト用ユーザーを作成（userID = 12345）
		ctx := context.Background()
		user, err := client.User.Create().
			SetEmail("test@example.com").
			SetPassword("dummy").
			SetIsAdmin(false).
			SetIsRoot(false).
			Save(ctx)
		assert.NoError(t, err)

		testSecret := "test_secret_key"
		err = os.Setenv("JWT_SECRET", testSecret)
		assert.NoError(t, err, "Setting JWT_SECRET should not produce an error")

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(middleware.AuthMiddleware())
		r.GET("/protected", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"userID": userID})
		})
		createdUserid := strconv.Itoa(user.ID)
		// 有効なJWTトークンを生成
		jwtGen := &utils.DefaultJWTGenerator{}
		token, _ := jwtGen.GenerateJWT(createdUserid)

		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"userID":`+createdUserid)
	})

	t.Run("TestAuthMiddleware_InvalidToken", func(t *testing.T) {
		// テスト用の JWT_SECRET を設定
		testSecret := "test_secret_key"
		err := os.Setenv("JWT_SECRET", testSecret)
		assert.NoError(t, err, "Setting JWT_SECRET should not produce an error")

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(middleware.AuthMiddleware())
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "You are authorized"})
		})

		// 無効なJWTトークンを使用
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"Invalid token"`)
	})

	t.Run("TestAuthMiddleware_NoToken", func(t *testing.T) {
		// テスト用の JWT_SECRET を設定
		testSecret := "test_secret_key"
		err := os.Setenv("JWT_SECRET", testSecret)
		assert.NoError(t, err, "Setting JWT_SECRET should not produce an error")

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(middleware.AuthMiddleware())
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "You are authorized"})
		})

		// トークンがないリクエスト
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"Authorization header required"`)
	})
}
