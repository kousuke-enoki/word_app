package word_test

import (
	"testing"
)

func TestWordList(_ *testing.T) {
	// gin.SetMode(gin.TestMode)

	// // Mockサービスの初期化
	// mockWordService := new(mocks.WordService)

	// // テスト用 JWT_SECRET 設定
	// testSecret := "test_secret_key"
	// err := os.Setenv("JWT_SECRET", testSecret)
	// assert.NoError(t, err)

	// // テスト用 JWT トークン作成
	// jwtGen := jwt.MyJWTGenerator{}
	// token, _ := jwtGen.GenerateJWT("1") // userID=1 を設定

	// t.Run("GetRegisteredWords_success", func(t *testing.T) {
	// 	// モックの戻り値設定
	// 	mockWordService := new(mocks.WordService)
	// 	mockWordService.On("GetRegisteredWords", mock.Anything, mock.Anything).Return(&models.WordListResponse{
	// 		Words: []models.Word{
	// 			{ID: 1, Name: "example"},
	// 		},
	// 		TotalPages: 1,
	// 	}, nil)
	// 	wordHandler := word.NewWordHandler(mockWordService)

	// 	// リクエスト準備
	// 	req, _ := http.NewRequest("GET", "/words?sortBy=register&order=asc&page=1&limit=10", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	// レスポンス記録
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
	// 	router.GET("/words", wordHandler.ListHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)
	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusOK, w.Code)

	// 	var resp models.WordListResponse
	// 	err = json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, 1, resp.TotalPages)
	// 	assert.Len(t, resp.Words, 1)
	// 	mockWordService.AssertExpectations(t)
	// })

	// t.Run("GetWords_success", func(t *testing.T) {
	// 	// モックの戻り値設定
	// 	mockWordService := new(mocks.WordService)
	// 	mockWordService.On("GetWords", mock.Anything, mock.Anything).Return(&models.WordListResponse{
	// 		Words: []models.Word{
	// 			{ID: 1, Name: "example"},
	// 		},
	// 		TotalPages: 1,
	// 	}, nil)
	// 	wordHandler := word.NewWordHandler(mockWordService)

	// 	// リクエスト準備
	// 	req, _ := http.NewRequest("GET", "/words?sortBy=name&order=asc&page=1&limit=10", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	// レスポンス記録
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
	// 	router.GET("/words", wordHandler.ListHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)
	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusOK, w.Code)

	// 	var resp models.WordListResponse
	// 	err = json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, 1, resp.TotalPages)
	// 	assert.Len(t, resp.Words, 1)
	// 	mockWordService.AssertExpectations(t)
	// })

	// t.Run("error: invalid query parameters", func(t *testing.T) {
	// 	// モックの戻り値設定
	// 	mockWordService := new(mocks.WordService)
	// 	mockWordService.On("GetWords", mock.Anything, mock.Anything).Return(nil, errors.New("invalid 'page' query parameter: must be a positive integer"))
	// 	wordHandler := word.NewWordHandler(mockWordService)

	// 	// リクエスト準備
	// 	req, _ := http.NewRequest("GET", "/words?sortBy=name&page=invalid", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	// レスポンス記録
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
	// 	router.GET("/words", wordHandler.ListHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)
	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 	// assert.Contains(t, ctx.Writer.Body.String(), "Invalid query parameters")

	// 	var resp map[string]string
	// 	err = json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, "invalid 'page' query parameter: must be a positive integer", resp["error"])
	// })

	// t.Run("error: userID not found in context", func(t *testing.T) {
	// 	wordHandler := word.NewWordHandler(mockWordService)
	// 	// リクエスト準備（Authorization ヘッダーなし）
	// 	req, _ := http.NewRequest("GET", "/words?sortBy=name&order=asc&page=1&limit=10", nil)

	// 	// レスポンス記録
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(func(c *gin.Context) {
	// 		// ユーザーIDをセットしないミドルウェア
	// 		c.Next()
	// 	})
	// 	router.GET("/words", wordHandler.ListHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)

	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	var resp map[string]string
	// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, "userID not found in context", resp["error"])
	// 	mockWordService.AssertExpectations(t)
	// })

	// t.Run("error: invalid userID type", func(t *testing.T) {
	// 	wordHandler := word.NewWordHandler(mockWordService)
	// 	// リクエスト準備
	// 	req, _ := http.NewRequest("GET", "/words?sortBy=name&order=asc&page=1&limit=10", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	// レスポンス記録
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(func(c *gin.Context) {
	// 		c.Set("userID", "invalid") // 不正な型の userID をセット
	// 		c.Next()
	// 	})
	// 	router.GET("/words", wordHandler.ListHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)

	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	var resp map[string]string
	// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, "invalid userID type", resp["error"])
	// 	mockWordService.AssertExpectations(t)
	// })

	// t.Run("error: service GetRegisteredWords fails", func(t *testing.T) {
	// 	// モックの戻り値設定
	// 	mockWordService := new(mocks.WordService)
	// 	mockWordService.On("GetRegisteredWords", mock.Anything, mock.Anything).Return(nil, errors.New("Service failure"))

	// 	wordHandler := word.NewWordHandler(mockWordService)

	// 	// リクエスト準備
	// 	req, _ := http.NewRequest("GET", "/words?sortBy=register&order=asc&page=1&limit=10", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	// レスポンス記録
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware())
	// 	router.GET("/words", wordHandler.ListHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)

	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// 	var resp map[string]string
	// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, "Service failure", resp["error"])
	// 	mockWordService.AssertExpectations(t)
	// })

	// t.Run("error: service GetWords fails", func(t *testing.T) {
	// 	// モックの戻り値設定
	// 	mockWordService := new(mocks.WordService)
	// 	mockWordService.On("GetWords", mock.Anything, mock.Anything).Return(nil, errors.New("Service failure"))

	// 	wordHandler := word.NewWordHandler(mockWordService)

	// 	// リクエスト準備
	// 	req, _ := http.NewRequest("GET", "/words?sortBy=name&order=asc&page=1&limit=10", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	// レスポンス記録
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware())
	// 	router.GET("/words", wordHandler.ListHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)

	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// 	var resp map[string]string
	// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, "Service failure", resp["error"])
	// 	mockWordService.AssertExpectations(t)
	// })
}
