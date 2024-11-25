// user/handler_test.go
package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/src/handlers/user"
	"word_app/backend/src/mocks"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// JWT トークン生成のモック関数
type MockJWTGenerator struct{}

func (m *MockJWTGenerator) GenerateJWT(userID string) (string, error) {
	return "mocked_jwt_token", nil
}

func TestSignUpHandler_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := new(mocks.UserClient)
	mockJWTGen := &MockJWTGenerator{}

	handler := user.NewUserHandler(mockClient, mockJWTGen)

	// 正常なリクエストデータ
	reqData := models.SignUpRequest{
		Email:    "test@example.com",
		Name:     "TestUser",
		Password: "Secure123!",
	}
	reqBody, _ := json.Marshal(reqData)

	mockClient.On("CreateUser", mock.Anything, reqData.Email, reqData.Name, mock.Anything).
		Return(&ent.User{ID: 1, Name: reqData.Name, Email: reqData.Email}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.SignUpHandler()(c)

	assert.Equal(t, http.StatusOK, w.Code)
	responseData := map[string]string{}
	err := json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "mocked_jwt_token", responseData["token"])

	mockClient.AssertExpectations(t)
}

func TestSignUpHandler_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := new(mocks.UserClient)
	mockJWTGen := &MockJWTGenerator{}

	handler := user.NewUserHandler(mockClient, mockJWTGen)

	// バリデーションエラーケース（短すぎる名前、無効なメール、簡単すぎるパスワード）
	testCases := []models.SignUpRequest{
		{Email: "invalid-email", Name: "T", Password: "simple"},
		{Email: "test@example.com", Name: "TestUser", Password: "short1"},
		{Email: "test@example.com", Name: "TestUser", Password: "NoSpecialChar1"},
	}

	for _, reqData := range testCases {
		reqBody, _ := json.Marshal(reqData)

		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.SignUpHandler()(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		responseData := map[string][]map[string]string{}
		err := json.Unmarshal(w.Body.Bytes(), &responseData)
		assert.NoError(t, err)

		// バリデーションエラーメッセージが返されていることを確認
		assert.Greater(t, len(responseData["errors"]), 0)
	}

	mockClient.AssertExpectations(t)
}
