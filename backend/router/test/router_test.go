package router_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	routerConfig "word_app/backend/router"
	"word_app/backend/src/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 環境変数設定
	_ = os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	// モックの作成
	mockUserHandler := new(mocks.UserHandler)
	mockWordHandler := new(mocks.WordHandler)

	// モックの動作設定
	mockUserHandler.On("SignUpHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "User signed up"})
		}
	})
	mockUserHandler.On("SignInHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "User signed in"})
		}
	})
	mockUserHandler.On("MyPageHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "this is MyPageHandler"})
		}
	})
	mockWordHandler.On("WordShowHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "WordShowHandler"})
		}
	})
	mockWordHandler.On("AllWordListHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "AllWordListHandler"})
		}
	})

	// Routerのセットアップ（テスト用にミドルウェアを無効化）
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Next()
	}) // シンプルなミドルウェアを挿入
	routerImpl := routerConfig.NewRouter(mockUserHandler, mockWordHandler)
	routerImpl.SetupRouter(router)

	// リクエスト送信
	req, _ := http.NewRequest("POST", "/users/sign_up", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 検証
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"User signed up"}`, w.Body.String())

	// モックの呼び出し検証
	mockUserHandler.AssertCalled(t, "SignUpHandler")
}
