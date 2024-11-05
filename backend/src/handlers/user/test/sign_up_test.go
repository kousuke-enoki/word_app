// user/handler_test.go
package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"word_app/backend/ent"
	"word_app/backend/src/handlers/user"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserClient は UserClient インターフェースをモック化したもの
type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) CreateUser(ctx context.Context, email, name, password string) (*ent.User, error) {
	args := m.Called(ctx, email, name, password)
	return args.Get(0).(*ent.User), args.Error(1)
}

func (m *MockUserClient) FindUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	return &ent.User{Email: email, Name: "Test User"}, nil
}

func (m *MockUserClient) FindUserByID(ctx context.Context, id int) (*ent.User, error) {
	return &ent.User{ID: id, Name: "Test User"}, nil
}

// JWT トークン生成のモック関数
type MockJWTGenerator struct{}

func (m *MockJWTGenerator) GenerateJWT(userID string) (string, error) {
	return "mocked_jwt_token", nil
}

func TestSignUpHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := new(MockUserClient)
	mockJWTGen := &MockJWTGenerator{} // モックの JWTGenerator を使用

	handler := user.NewUserHandler(mockClient, mockJWTGen) // モックの JWTGenerator を注入

	reqData := models.SignUpRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "securepassword",
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
