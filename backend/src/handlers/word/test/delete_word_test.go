package word_test

import (
	"testing"
)

func TestDeleteWordHandler(t *testing.T) {
	// gin.SetMode(gin.TestMode)

	// // テスト用 JWT_SECRET 設定
	// testSecret := "test_secret_key"
	// err := os.Setenv("JWT_SECRET", testSecret)
	// assert.NoError(t, err)

	// // テスト用 JWT トークン作成
	// jwtGen := &jwt.MyJWTGenerator{}
	// token, _ := jwtGen.GenerateJWT("1") // userID=1 を設定

	// t.Run("DeleteWord_success", func(t *testing.T) {
	// 	mockService := new(mocks.WordService)
	// 	mockService.On("DeleteWord", mock.Anything, mock.Anything).Return(&models.DeleteWordResponse{
	// 		Name: "test", Message: "word delete complete"}, nil)
	// 	wordHandler := word.NewWordHandler(mockService)

	// 	// リクエスト準備
	// 	req, _ := http.NewRequest("DELETE", "/words/1", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	// レスポンス記録
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware()) // テスト用ミドルウェア
	// 	router.DELETE("/words/:id", wordHandler.DeleteWordHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)
	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusOK, w.Code)

	// 	var resp models.DeleteWordResponse
	// 	err = json.Unmarshal(w.Body.Bytes(), &resp)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, "test", resp.Name)
	// 	mockService.AssertExpectations(t)
	// })

	// t.Run("InvalidWordID_BadRequest", func(t *testing.T) {
	// 	mockService := new(mocks.WordService)
	// 	wordHandler := word.NewWordHandler(mockService)

	// 	// 無効なWord IDをリクエスト
	// 	req, _ := http.NewRequest("DELETE", "/words/invalid_id", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware())
	// 	router.DELETE("/words/:id", wordHandler.DeleteWordHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)

	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	assert.Contains(t, w.Body.String(), "Invalid word ID")
	// })

	// t.Run("UnauthorizedUser_MissingUserID", func(t *testing.T) {
	// 	mockService := new(mocks.WordService)
	// 	wordHandler := word.NewWordHandler(mockService)

	// 	// ヘッダーなしでリクエスト
	// 	req, _ := http.NewRequest("DELETE", "/words/1", nil)

	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddlewareWithoutUser()) // userIDをセットしないモック
	// 	router.DELETE("/words/:id", wordHandler.DeleteWordHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)

	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	assert.Contains(t, w.Body.String(), "unauthorized: userID not found in context")
	// })

	// t.Run("InvalidUserIDType_BadRequest", func(t *testing.T) {
	// 	mockService := new(mocks.WordService)
	// 	wordHandler := word.NewWordHandler(mockService)

	// 	// userIDに不正な型を設定
	// 	req, _ := http.NewRequest("DELETE", "/words/1", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddlewareWithInvalidUserType()) // 不正なuserID型をセットするモック
	// 	router.DELETE("/words/:id", wordHandler.DeleteWordHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)

	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	assert.Contains(t, w.Body.String(), "invalid userID type")
	// })

	// t.Run("InternalServerError_GetWordDetails", func(t *testing.T) {
	// 	mockService := new(mocks.WordService)
	// 	mockService.On("DeleteWord", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
	// 	wordHandler := word.NewWordHandler(mockService)

	// 	// 正常なリクエスト
	// 	req, _ := http.NewRequest("DELETE", "/words/1", nil)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.Use(mocks.MockAuthMiddleware())
	// 	router.DELETE("/words/:id", wordHandler.DeleteWordHandler())

	// 	// テスト実行
	// 	router.ServeHTTP(w, req)

	// 	// レスポンス検証
	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// 	assert.Contains(t, w.Body.String(), "database error")
	// })
}
