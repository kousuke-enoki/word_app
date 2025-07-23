package router

import (
	"testing"
)

func RouterTest(_ *testing.T) {
	// t.Run("TestSetupRouter", func(t *testing.T) {
	// 	gin.SetMode(gin.TestMode)

	// 	// 環境変数設定
	// 	_ = os.Setenv("JWT_SECRET", "test_secret")
	// 	defer os.Unsetenv("JWT_SECRET")

	// 	// モックの作成
	// 	mockUserHandler := new(mocks.UserHandler)
	// 	mockSettingHandler := new(mocks.SettingHandler)
	// 	mockWordHandler := new(mocks.WordHandler)
	// 	mockAuthHandler := new(mocks.AuthHandler)

	// 	// モックの動作設定
	// 	mockUserHandler.On("SignUpHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "User signed up"})
	// 		}
	// 	})
	// 	mockUserHandler.On("SignInHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "User signed in"})
	// 		}
	// 	})
	// 	mockAuthHandler.On("AuthCheckHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "this is AuthCheckHandler"})
	// 		}
	// 	})
	// 	mockUserHandler.On("MyPageHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "this is MyPageHandler"})
	// 		}
	// 	})
	// 	mockWordHandler.On("CreateHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "this is CreateHandler"})
	// 		}
	// 	})
	// 	mockWordHandler.On("UpdateHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "this is UpdateHandler"})
	// 		}
	// 	})
	// 	mockWordHandler.On("DeleteHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "this is DeleteHandler"})
	// 		}
	// 	})
	// 	mockWordHandler.On("ListHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "ListHandler"})
	// 		}
	// 	})
	// 	mockWordHandler.On("ShowHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "ShowHandler"})
	// 		}
	// 	})
	// 	mockWordHandler.On("ShowHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "ShowHandler"})
	// 		}
	// 	})
	// 	mockWordHandler.On("RegisterHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "RegisterHandler"})
	// 		}
	// 	})
	// 	mockWordHandler.On("SaveMemoHandler").Return(func() gin.HandlerFunc {
	// 		return func(c *gin.Context) {
	// 			c.JSON(http.StatusOK, gin.H{"message": "SaveMemoHandler"})
	// 		}
	// 	})

	// 	// Routerのセットアップ（テスト用にミドルウェアを無効化）
	// 	router := gin.New()
	// 	router.Use(func(c *gin.Context) {
	// 		c.Next()
	// 	}) // シンプルなミドルウェアを挿入
	// 	routerImpl := NewRouter(
	// 		mockAuthHandler,
	// 		mockUserHandler,
	// 		mockSettingHandler,
	// 		mockWordHandler,
	// 	)
	// 	routerImpl.SetupRouter(router)

	// 	// リクエスト送信
	// 	req, _ := http.NewRequest("POST", "/users/sign_up", nil)
	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)

	// 	// 検証
	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	assert.JSONEq(t, `{"message":"User signed up"}`, w.Body.String())

	// 	// モックの呼び出し検証
	// 	mockUserHandler.AssertCalled(t, "SignUpHandler")
	// })

	// t.Run("TestCORSMiddleware", func(t *testing.T) {
	// 	router := gin.New()
	// 	router.Use(CORSMiddleware())
	// 	router.OPTIONS("/test", func(c *gin.Context) {
	// 		c.JSON(http.StatusNoContent, nil)
	// 	})

	// 	w := httptest.NewRecorder()
	// 	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	// 	router.ServeHTTP(w, req)

	// 	if w.Code != http.StatusNoContent {
	// 		t.Errorf("Expected status code 204 for OPTIONS request, got %d", w.Code)
	// 	}
	// })

	// t.Run("TestRequestLoggerMiddleware", func(t *testing.T) {
	// 	router := gin.New()
	// 	router.Use(requestLoggerMiddleware())
	// 	router.GET("/test", func(c *gin.Context) {
	// 		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	// 	})

	// 	w := httptest.NewRecorder()
	// 	req, _ := http.NewRequest("GET", "/test", nil)
	// 	router.ServeHTTP(w, req)

	// 	if w.Code != http.StatusOK {
	// 		t.Errorf("Expected status code 200, got %d", w.Code)
	// 	}
	// })
}
