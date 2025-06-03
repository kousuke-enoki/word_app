package word_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"word_app/backend/src/handlers/word"
	"word_app/backend/src/infrastructure/jwt"
	"word_app/backend/src/mocks"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveMemoHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// mockWordService := new(mocks.WordService)

	// テスト用 JWT_SECRET 設定
	testSecret := "test_secret_key"
	err := os.Setenv("JWT_SECRET", testSecret)
	assert.NoError(t, err)

	// テスト用 JWT トークン作成
	jwtGen := jwt.MyJWTGenerator{}
	token, _ := jwtGen.GenerateJWT("1") // userID=1 を設定

	t.Run("SaveMemo_success", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		mockWordService.On("SaveMemo", mock.Anything, mock.Anything).Return(&models.SaveMemoResponse{
			Name: "test", Memo: "Memo test", Message: "RegisteredWord updated",
		}, nil)
		wordHandler := word.NewWordHandler(mockWordService)

		// 正常なリクエストデータ
		reqData := models.SaveMemoRequest{
			WordID: 1,
			Memo:   "Memo test",
		}

		reqBody, _ := json.Marshal(reqData)
		// リクエスト準備
		req, _ := http.NewRequest(http.MethodPost, "/words/memo", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
		router.POST("/words/memo", wordHandler.SaveMemoHandler())

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		assert.Equal(t, http.StatusOK, w.Code)

		var resp models.SaveMemoResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "test", resp.Name)
		mockWordService.AssertExpectations(t)
	})
	t.Run("SaveMemo_invalidJSON", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		wordHandler := word.NewWordHandler(mockWordService)

		// 不正なJSONフォーマット
		reqBody := []byte(`{"invalidJson":}`)

		// リクエスト準備
		req, _ := http.NewRequest(http.MethodPost, "/words/memo", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
		router.POST("/words/memo", wordHandler.SaveMemoHandler())

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid JSON format")
	})

	t.Run("SaveMemo_missingRequestBody", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		wordHandler := word.NewWordHandler(mockWordService)

		// リクエストボディなし
		req, _ := http.NewRequest(http.MethodPost, "/words/memo", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
		router.POST("/words/memo", wordHandler.SaveMemoHandler())

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "request body is missing")
	})

	t.Run("SaveMemo_validationErrors", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		wordHandler := word.NewWordHandler(mockWordService)

		// バリデーションエラーを引き起こすデータ
		reqData := models.SaveMemoRequest{
			WordID: 0, // WordIDが無効
			Memo:   "",
		}

		reqBody, _ := json.Marshal(reqData)
		// リクエスト準備
		req, _ := http.NewRequest(http.MethodPost, "/words/memo", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
		router.POST("/words/memo", wordHandler.SaveMemoHandler())

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("SaveMemo_internalServerError", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		mockWordService.On("SaveMemo", mock.Anything, mock.Anything).Return(nil, errors.New("internal server error"))
		wordHandler := word.NewWordHandler(mockWordService)

		// 正常なリクエストデータ
		reqData := models.SaveMemoRequest{
			WordID: 1,
			Memo:   "Valid Memo",
		}

		reqBody, _ := json.Marshal(reqData)
		// リクエスト準備
		req, _ := http.NewRequest(http.MethodPost, "/words/memo", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
		router.POST("/words/memo", wordHandler.SaveMemoHandler())

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

	t.Run("SaveMemo_missingUserIDInContext", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		wordHandler := word.NewWordHandler(mockWordService)

		// 正常なリクエストデータ
		reqData := models.SaveMemoRequest{
			WordID: 1,
			Memo:   "Valid Memo",
		}

		reqBody, _ := json.Marshal(reqData)
		// リクエスト準備
		req, _ := http.NewRequest(http.MethodPost, "/words/memo", bytes.NewBuffer(reqBody))
		// ヘッダーに認証情報を付加しない

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/words/memo", wordHandler.SaveMemoHandler()) // ミドルウェアを適用しない

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized: userID not found in context")
	})
}
