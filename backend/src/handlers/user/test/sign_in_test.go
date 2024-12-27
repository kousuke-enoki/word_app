package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"golang.org/x/crypto/bcrypt"
)

func TestSignInHandler(t *testing.T) {
	t.Run("TestSignInHandler_ValidRequest", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockClient := new(mocks.UserClient)
		mockJWTGen := new(mocks.MockJwtGenerator)

		handler := user.NewUserHandler(mockClient, mockJWTGen)

		// 正常なリクエストデータ
		reqData := models.SignInRequest{
			Email:    "test@example.com",
			Password: "Secure123!",
		}
		reqBody, _ := json.Marshal(reqData)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Secure123!"), bcrypt.DefaultCost)
		mockClient.On("FindUserByEmail", mock.Anything, reqData.Email).
			Return(&ent.User{ID: 1, Email: reqData.Email, Password: string(hashedPassword)}, nil)
		mockJWTGen.On("GenerateJWT", "1").Return("mocked_jwt_token", nil)

		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.SignInHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)

		responseData := map[string]string{}
		err := json.Unmarshal(w.Body.Bytes(), &responseData)
		assert.NoError(t, err)
		assert.Equal(t, "mocked_jwt_token", responseData["token"])
		assert.Equal(t, "Authentication successful", responseData["message"])

		mockClient.AssertExpectations(t)
		mockJWTGen.AssertExpectations(t)
	})

	t.Run("TestSignInHandler_InvalidCredentials", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockClient := new(mocks.UserClient)
		mockJWTGen := &mocks.MockJwtGenerator{}
		handler := user.NewUserHandler(mockClient, mockJWTGen)

		// 無効なリクエストデータ
		password := "InvalidSecure123!"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Secure123!"), bcrypt.DefaultCost)
		signInUser := &ent.User{
			ID:       1,
			Email:    "test@example.com",
			Name:     "Test User",
			Password: string(hashedPassword),
		}

		// モックの設定
		mockClient.On("FindUserByEmail", mock.Anything, "test@example.com").Return(signInUser, nil)

		// リクエスト作成
		reqData := models.SignInRequest{
			Email:    "test@example.com",
			Password: password,
		}
		reqBody, _ := json.Marshal(reqData)
		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// レスポンス作成
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// ハンドラ呼び出し
		handler.SignInHandler()(c)

		// 検証
		assert.Equal(t, http.StatusBadRequest, w.Code)
		responseData := map[string]string{}
		err := json.Unmarshal(w.Body.Bytes(), &responseData)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request", responseData["error"])

		mockClient.AssertExpectations(t)
		mockJWTGen.AssertExpectations(t)
	})

	t.Run("TestSignInHandler_TokenGenerationError", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockClient := new(mocks.UserClient)
		mockJWTGen := &mocks.MockJwtGenerator{}
		handler := user.NewUserHandler(mockClient, mockJWTGen)

		// 正常なユーザーデータとトークン生成エラーのモック設定
		password := "Secure123!"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		signInUser := &ent.User{
			ID:       1,
			Email:    "test@example.com",
			Name:     "Test User",
			Password: string(hashedPassword),
		}

		mockClient.On("FindUserByEmail", mock.Anything, "test@example.com").Return(signInUser, nil)
		mockJWTGen.On("GenerateJWT", "1").Return("", fmt.Errorf("token generation error"))

		// リクエスト作成
		reqData := models.SignInRequest{
			Email:    "test@example.com",
			Password: password,
		}
		reqBody, _ := json.Marshal(reqData)
		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// レスポンス作成
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// ハンドラ呼び出し
		handler.SignInHandler()(c)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		responseData := map[string]string{}
		err := json.Unmarshal(w.Body.Bytes(), &responseData)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to generate token", responseData["error"])

		mockClient.AssertExpectations(t)
		mockJWTGen.AssertExpectations(t)
	})
}
