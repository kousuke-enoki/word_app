package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/src/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRootHandler(t *testing.T) {
	// Ginのテストモードを有効化
	gin.SetMode(gin.TestMode)

	// ルーターとハンドラー設定
	router := gin.Default()
	router.GET("/", handlers.RootHandler)

	// テスト用のリクエストを作成
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// ハンドラー実行
	router.ServeHTTP(w, req)

	// ステータスコードの検証
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスの検証
	expectedResponse := gin.H{
		"message":      "Redirect to root",
		"redirect_url": "/",
	}
	var actualResponse gin.H
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}
