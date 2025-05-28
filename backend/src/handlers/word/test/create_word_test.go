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
	"word_app/backend/src/infrastructure/auth"
	"word_app/backend/src/mocks"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateWordHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// mockWordService := new(mocks.WordService)

	// テスト用 JWT_SECRET 設定
	testSecret := "test_secret_key"
	err := os.Setenv("JWT_SECRET", testSecret)
	assert.NoError(t, err)

	// テスト用 JWT トークン作成
	jwtGen := &auth.DefaultJWTGenerator{}
	token, _ := jwtGen.GenerateJWT("1") // userID=1 を設定

	t.Run("CreateWord_success", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		mockWordService.On("CreateWord", mock.Anything, mock.Anything).Return(&models.CreateWordResponse{
			ID: 1, Name: "test", Message: "RegisteredWord updated",
		}, nil)
		wordHandler := word.NewWordHandler(mockWordService)

		// 正常なリクエストデータ
		japaneseMean := []models.JapaneseMean{
			{
				Name: "テスト",
			},
		}
		wordInfos := []models.WordInfo{
			{
				PartOfSpeechID: 1,
				JapaneseMeans:  japaneseMean,
			},
		}
		reqData := models.CreateWordRequest{
			Name:      "test",
			WordInfos: wordInfos,
		}

		reqBody, _ := json.Marshal(reqData)
		// リクエスト準備
		req, _ := http.NewRequest(http.MethodPost, "/words/new", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
		router.POST("/words/new", wordHandler.CreateWordHandler())

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		assert.Equal(t, http.StatusOK, w.Code)

		var resp models.CreateWordResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "test", resp.Name)
		mockWordService.AssertExpectations(t)
	})
	t.Run("CreateWord_fail", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		mockWordService.On("CreateWord", mock.Anything, mock.Anything).Return(nil, errors.New("word with the same name already exists"))
		wordHandler := word.NewWordHandler(mockWordService)

		// 正常なリクエストデータ
		japaneseMean := []models.JapaneseMean{
			{
				Name: "テスト",
			},
		}
		wordInfos := []models.WordInfo{
			{
				PartOfSpeechID: 1,
				JapaneseMeans:  japaneseMean,
			},
		}
		reqData := models.CreateWordRequest{
			Name:      "test",
			WordInfos: wordInfos,
		}

		reqBody, _ := json.Marshal(reqData)
		// リクエスト準備
		req, _ := http.NewRequest(http.MethodPost, "/words/new", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
		router.POST("/words/new", wordHandler.CreateWordHandler())

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		var resp gin.H

		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		// err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "Failed to create word", resp["error"])
	})
	t.Run("CreateWord_invalidJSON", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		wordHandler := word.NewWordHandler(mockWordService)

		// 無効なJSONデータ
		invalidJSON := `{"Name": "test", "WordInfos": "wordInfos"}`

		req, _ := http.NewRequest(http.MethodPost, "/words/new", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware())
		router.POST("/words/new", wordHandler.CreateWordHandler())

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp gin.H
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "json: cannot unmarshal string into Go struct field CreateWordRequest.wordInfos of type []models.WordInfo", resp["error"])
	})

	t.Run("CreateWord_missingRequestBody", func(t *testing.T) {
		mockWordService := new(mocks.WordService)
		wordHandler := word.NewWordHandler(mockWordService)

		// リクエストボディなし
		req, _ := http.NewRequest(http.MethodPost, "/words/new", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		// レスポンス記録
		w := httptest.NewRecorder()
		router := gin.Default()
		router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
		router.POST("/words/new", wordHandler.CreateWordHandler())

		// テスト実行
		router.ServeHTTP(w, req)
		// レスポンス検証
		var resp gin.H

		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		// err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "invalid request", resp["error"])
	})

	// t.Run("CreateWord_validateError", func(t *testing.T) {
	// 	mockWordService := new(mocks.WordService)
	// 	mockWordService.On("CreateWord", mock.Anything, mock.Anything).Return(&models.CreateWordResponse{
	// 		ID: 1, Name: "test", Message: "RegisteredWord updated",
	// 	}, nil)
	// 	wordHandler := word.NewWordHandler(mockWordService)

	// 	// 正常なリクエストデータ
	// 	japaneseMean := []models.JapaneseMean{
	// 		{
	// 			Name: "test",
	// 		},
	// 	}
	// 	wordInfos := []models.WordInfo{
	// 		{
	// 			PartOfSpeechID: 1,
	// 			JapaneseMeans:  japaneseMean,
	// 		},
	// 	}
	// 	reqData := models.CreateWordRequest{
	// 		Name:      "テスト",
	// 		WordInfos: wordInfos,
	// 	}

	// 	reqBody, _ := json.Marshal(reqData)
	// 	// リクエスト準備
	// 	req, _ := http.NewRequest(http.MethodPost, "/words/new", bytes.NewBuffer(reqBody))
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	// レスポンス記録
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
	// 	router.POST("/words/new", wordHandler.CreateWordHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)
	// 	// レスポンス検証
	// 	var resp gin.H
	// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	// err := json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	assert.Equal(t, "invalid request", resp["error"])
	// })
}
