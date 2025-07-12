package word_test

import (
	"testing"
)

func TestUpdateWordHandler(t *testing.T) {
	// gin.SetMode(gin.TestMode)
	// // mockWordService := new(mocks.WordService)

	// // テスト用 JWT_SECRET 設定
	// testSecret := "test_secret_key"
	// err := os.Setenv("JWT_SECRET", testSecret)
	// assert.NoError(t, err)

	// // テスト用 JWT トークン作成
	// jwtGen := jwt.MyJWTGenerator{}
	// token, _ := jwtGen.GenerateJWT("1") // userID=1 を設定

	// t.Run("UpdateWord_success", func(t *testing.T) {
	// 	mockWordService := new(mocks.WordService)
	// 	mockWordService.On("UpdateWord", mock.Anything, mock.Anything).Return(&models.UpdateWordResponse{
	// 		ID: 1, Name: "test", Message: "RegisteredWord updated",
	// 	}, nil)
	// 	wordHandler := word.NewWordHandler(mockWordService)

	// 	// 正常なリクエストデータ
	// 	japaneseMean := []models.JapaneseMean{
	// 		{
	// 			Name: "テスト",
	// 		},
	// 	}
	// 	wordInfos := []models.WordInfo{
	// 		{
	// 			PartOfSpeechID: 1,
	// 			JapaneseMeans:  japaneseMean,
	// 		},
	// 	}
	// 	reqData := models.UpdateWordRequest{
	// 		Name:      "test",
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
	// 	router.POST("/words/new", wordHandler.UpdateWordHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)
	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusOK, w.Code)

	// 	var resp models.UpdateWordResponse
	// 	err = json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, "test", resp.Name)
	// 	mockWordService.AssertExpectations(t)
	// })
	// t.Run("UpdateWord_fail", func(t *testing.T) {
	// 	mockWordService := new(mocks.WordService)
	// 	mockWordService.On("UpdateWord", mock.Anything, mock.Anything).Return(nil, errors.New("word with the same name already exists"))
	// 	wordHandler := word.NewWordHandler(mockWordService)

	// 	// 正常なリクエストデータ
	// 	japaneseMean := []models.JapaneseMean{
	// 		{
	// 			Name: "テスト",
	// 		},
	// 	}
	// 	wordInfos := []models.WordInfo{
	// 		{
	// 			PartOfSpeechID: 1,
	// 			JapaneseMeans:  japaneseMean,
	// 		},
	// 	}
	// 	reqData := models.UpdateWordRequest{
	// 		Name:      "test",
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
	// 	router.POST("/words/new", wordHandler.UpdateWordHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)
	// 	// レスポンス検証
	// 	var resp gin.H

	// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	// err := json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// 	assert.Equal(t, "Failed to update word", resp["error"])
	// })
	// t.Run("UpdateWord_invalidJSON", func(t *testing.T) {
	// 	mockWordService := new(mocks.WordService)
	// 	wordHandler := word.NewWordHandler(mockWordService)

	// 	// 無効なJSONデータ
	// 	invalidJSON := `{"Name": "test", "WordInfos": "wordInfos"}`

	// 	req, _ := http.NewRequest(http.MethodPost, "/words/new", bytes.NewBufferString(invalidJSON))
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware())
	// 	router.POST("/words/new", wordHandler.UpdateWordHandler())

	// 	router.ServeHTTP(w, req)
	// 	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 	var resp gin.H
	// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, "json: cannot unmarshal string into Go struct field UpdateWordRequest.wordInfos of type []models.WordInfo", resp["error"])
	// })

	// t.Run("UpdateWord_missingRequestBody", func(t *testing.T) {
	// 	mockWordService := new(mocks.WordService)
	// 	wordHandler := word.NewWordHandler(mockWordService)

	// 	// リクエストボディなし
	// 	req, _ := http.NewRequest(http.MethodPost, "/words/new", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	// レスポンス記録
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
	// 	router.POST("/words/new", wordHandler.UpdateWordHandler())

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

	// // t.Run("UpdateWord_validateError", func(t *testing.T) {
	// // 	mockWordService := new(mocks.WordService)
	// // 	mockWordService.On("UpdateWord", mock.Anything, mock.Anything).Return(&models.UpdateWordResponse{
	// // 		ID: 1, Name: "test", Message: "RegisteredWord updated",
	// // 	}, nil)
	// // 	wordHandler := word.NewWordHandler(mockWordService)

	// // 	// 正常なリクエストデータ
	// // 	japaneseMean := []models.JapaneseMean{
	// // 		{
	// // 			Name: "test",
	// // 		},
	// // 	}
	// // 	wordInfos := []models.WordInfo{
	// // 		{
	// // 			PartOfSpeechID: 1,
	// // 			JapaneseMeans:  japaneseMean,
	// // 		},
	// // 	}
	// // 	reqData := models.UpdateWordRequest{
	// // 		Name:      "テスト",
	// // 		WordInfos: wordInfos,
	// // 	}

	// // 	reqBody, _ := json.Marshal(reqData)
	// // 	// リクエスト準備
	// // 	req, _ := http.NewRequest(http.MethodPost, "/words/new", bytes.NewBuffer(reqBody))
	// // 	req.Header.Set("Authorization", "Bearer "+token)

	// // 	// レスポンス記録
	// // 	w := httptest.NewRecorder()
	// // 	router := gin.Default()
	// // 	router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
	// // 	router.POST("/words/new", wordHandler.UpdateWordHandler())

	// // 	// テスト実行
	// // 	router.ServeHTTP(w, req)
	// // 	// レスポンス検証
	// // 	var resp gin.H
	// // 	err := json.Unmarshal(w.Body.Bytes(), &resp)
	// // 	assert.NoError(t, err)
	// // 	// err := json.Unmarshal(w.Body.Bytes(), &resp)
	// // 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// // 	assert.Equal(t, "invalid request", resp["error"])
	// // })
}
