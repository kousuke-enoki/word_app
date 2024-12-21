package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
	mockAuthHandler := new(mocks.AuthHandler)

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
	mockAuthHandler.On("AuthCheckHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "this is AuthCheckHandler"})
		}
	})
	mockUserHandler.On("MyPageHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "this is MyPageHandler"})
		}
	})
	mockWordHandler.On("CreateWordHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "this is CreateWordHandler"})
		}
	})
	mockWordHandler.On("UpdateWordHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "this is UpdateWordHandler"})
		}
	})
	mockWordHandler.On("DeleteWordHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "this is DeleteWordHandler"})
		}
	})
	mockWordHandler.On("AllWordListHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "AllWordListHandler"})
		}
	})
	mockWordHandler.On("WordShowHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "WordShowHandler"})
		}
	})
	mockWordHandler.On("WordShowHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "WordShowHandler"})
		}
	})
	mockWordHandler.On("RegisterWordHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "RegisterWordHandler"})
		}
	})
	mockWordHandler.On("SaveMemoHandler").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "SaveMemoHandler"})
		}
	})

	// Routerのセットアップ（テスト用にミドルウェアを無効化）
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Next()
	}) // シンプルなミドルウェアを挿入
	routerImpl := NewRouter(mockAuthHandler, mockUserHandler, mockWordHandler)
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

func TestCORSMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(CORSMiddleware())
	router.OPTIONS("/test", func(c *gin.Context) {
		c.JSON(http.StatusNoContent, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code 204 for OPTIONS request, got %d", w.Code)
	}
}

func TestRequestLoggerMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(requestLoggerMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}
