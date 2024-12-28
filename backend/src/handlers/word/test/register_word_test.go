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
	"word_app/backend/src/mocks"
	"word_app/backend/src/models"
	"word_app/backend/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterWordHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// mockWordService := new(mocks.WordService)

	// テスト用 JWT_SECRET 設定
	testSecret := "test_secret_key"
	err := os.Setenv("JWT_SECRET", testSecret)
	assert.NoError(t, err)

	// テスト用 JWT トークン作成
	jwtGen := &utils.DefaultJWTGenerator{}
	token, _ := jwtGen.GenerateJWT("1") // userID=1 を設定

	t.Run("RegisterWord_success", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		mockWordService.On("RegisterWords", mock.Anything, mock.Anything).Return(&models.RegisterWordResponse{
			Name: "test", IsRegistered: true, RegistrationCount: 1, Message: "RegisteredWord updated",
		}, nil)
		wordHandler := word.NewWordHandler(mockWordService)

		// 正常なリクエストデータ
		reqData := models.RegisterWordRequest{
			WordID:       1,
			IsRegistered: false,
		}

		reqBody, _ := json.Marshal(reqData)
		// リクエスト準備
		req, _ := http.NewRequest(http.MethodPost, "/words/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
		router.POST("/words/register", wordHandler.RegisterWordHandler())

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		assert.Equal(t, http.StatusOK, w.Code)

		var resp models.RegisterWordResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "test", resp.Name)
		mockWordService.AssertExpectations(t)
	})

	t.Run("UnRegisterWord_success", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		mockWordService.On("RegisterWords", mock.Anything, mock.Anything).Return(&models.RegisterWordResponse{
			Name: "test", IsRegistered: false, RegistrationCount: 1, Message: "RegisteredWord updated",
		}, nil)
		wordHandler := word.NewWordHandler(mockWordService)

		// 正常なリクエストデータ
		reqData := models.RegisterWordRequest{
			WordID:       1,
			IsRegistered: true,
		}

		reqBody, _ := json.Marshal(reqData)
		// リクエスト準備
		req, _ := http.NewRequest(http.MethodPost, "/words/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
		router.POST("/words/register", wordHandler.RegisterWordHandler())

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		assert.Equal(t, http.StatusOK, w.Code)

		var resp models.RegisterWordResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "test", resp.Name)
		mockWordService.AssertExpectations(t)
	})
	t.Run("RegisterWord_invalidJSON", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		wordHandler := word.NewWordHandler(mockWordService)

		// 無効なJSONデータ
		invalidJSON := `{"WordID": 1, "IsRegistered": "invalid_bool"}`

		req, _ := http.NewRequest(http.MethodPost, "/words/register", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware())
		router.POST("/words/register", wordHandler.RegisterWordHandler())

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp gin.H
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "invalid JSON format: json: cannot unmarshal string into Go struct field RegisterWordRequest.isRegistered of type bool", resp["error"])
	})

	t.Run("RegisterWord_missingRequestBody", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		wordHandler := word.NewWordHandler(mockWordService)

		req, _ := http.NewRequest(http.MethodPost, "/words/register", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware())
		router.POST("/words/register", wordHandler.RegisterWordHandler())

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp gin.H
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "request body is missing", resp["error"])
	})

	t.Run("RegisterWord_serviceError", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		mockWordService.On("RegisterWords", mock.Anything, mock.Anything).Return(nil, errors.New("internal service error"))
		wordHandler := word.NewWordHandler(mockWordService)

		reqData := models.RegisterWordRequest{
			WordID:       1,
			IsRegistered: false,
		}

		reqBody, _ := json.Marshal(reqData)
		req, _ := http.NewRequest(http.MethodPost, "/words/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware())
		router.POST("/words/register", wordHandler.RegisterWordHandler())

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp gin.H
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "internal service error", resp["error"])
	})

	t.Run("RegisterWord_missingUserID", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		wordHandler := word.NewWordHandler(mockWordService)

		reqData := models.RegisterWordRequest{
			WordID:       1,
			IsRegistered: false,
		}

		reqBody, _ := json.Marshal(reqData)
		req, _ := http.NewRequest(http.MethodPost, "/words/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/words/register", wordHandler.RegisterWordHandler()) // Middlewareなし

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp gin.H
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "unauthorized: userID not found in context", resp["error"])
	})
}
