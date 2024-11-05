package user_test

import (
	"testing"
)

func TestSignUpHandler(t *testing.T) {
	// // モックを生成
	// mockClient := new(mocks.UserClient)

	// // フェイクデータ生成
	// reqBody := models.SignUpRequest{}
	// _ = faker.FakeData(&reqBody)

	// // モックの期待値を設定
	// hashedPassword := "$2a$10$..." // ハッシュ化済みのパスワード
	// mockClient.On("CreateUser", mock.Anything, reqBody.Email, reqBody.Name, hashedPassword).Return(&ent.User{ID: 1, Email: reqBody.Email, Name: reqBody.Name, Password: hashedPassword}, nil)

	// // Ginコンテキストの準備
	// gin.SetMode(gin.TestMode)
	// router := gin.Default()
	// handlers := handlers.NewUserHandler(mockClient)
	// router.POST("/users/sign_up", handlers.SignUpHandler())

	// // リクエストをJSONに変換
	// body, _ := json.Marshal(reqBody)
	// req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	// req.Header.Set("Content-Type", "application/json")
	// w := httptest.NewRecorder()

	// // テストの実行
	// router.ServeHTTP(w, req)

	// // 検証
	// assert.Equal(t, http.StatusOK, w.Code)
	// mockClient.AssertExpectations(t)
}
