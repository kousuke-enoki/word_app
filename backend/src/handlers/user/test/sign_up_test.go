package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/src/handlers/user"
	"word_app/backend/src/mocks"
	"word_app/backend/src/models"

	user_mocks "word_app/backend/src/mocks/http/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUpHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// モックの初期化
	mockClient := new(user_mocks.MockUsecase)
	mockJWTGen := &mocks.MockJwtGenerator{}
	handler := user.NewHandler(mockClient, mockJWTGen)

	t.Run("error: request body is nil", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/sign_up", nil)
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/sign_up", handler.SignUpHandler())

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "request body is nil", resp["error"])
		mockClient.AssertExpectations(t)
	})

	t.Run("error: invalid JSON request", func(t *testing.T) {
		reqBody := `{"invalidJson":`
		req, _ := http.NewRequest("POST", "/sign_up", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/sign_up", handler.SignUpHandler())

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp["error"], "invalid request")
		mockClient.AssertExpectations(t)
	})

	t.Run("error: validation errors", func(t *testing.T) {
		reqBody := `{"email": "", "name": "", "password": ""}`
		req, _ := http.NewRequest("POST", "/sign_up", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/sign_up", handler.SignUpHandler())

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "invalid request: Key: 'SignUpRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'SignUpRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'SignUpRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag", resp["error"])
		mockClient.AssertExpectations(t)
	})

	t.Run("error: validate input fields", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		userHandler := user.NewHandler(mockClient, mockJWTGen)

		// モックの設定
		mockClient.On("Create", mock.Anything, "test@example.com", "te", "pass").
			Return(nil, gin.H{"error": "ErrDuplicateEmail"})

		// リクエストの設定
		req, _ := http.NewRequest("POST", "/sign_up", strings.NewReader(`{"email":"test@example.com","name":"te","password":"pass"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Ginルーターの設定
		router := gin.Default()
		router.POST("/sign_up", userHandler.SignUpHandler())

		// テストリクエストの送信
		router.ServeHTTP(w, req)

		// ステータスコードの確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// レスポンスの解析
		var resp struct {
			Errors []struct {
				Field   string `json:"field"`
				Message string `json:"message"`
			} `json:"errors"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// エラーメッセージの確認
		assert.Len(t, resp.Errors, 3)

		// エラー内容の確認
		assert.Equal(t, "name", resp.Errors[0].Field)
		assert.Equal(t, "name must be between 3 and 20 characters", resp.Errors[0].Message)

		assert.Equal(t, "password", resp.Errors[1].Field)
		assert.Equal(t, "password must be between 8 and 30 characters", resp.Errors[1].Message)

		assert.Equal(t, "password", resp.Errors[2].Field)
		assert.Equal(t, "password must include at least one uppercase letter, one lowercase letter, one number, and one special character", resp.Errors[2].Message)
	})

	t.Run("error: database failure", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		userHandler := user.NewHandler(mockClient, mockJWTGen)
		// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123$"), bcrypt.DefaultCost)

		mockClient.On("Create", mock.Anything, "test@example.com", "test", mock.Anything).
			Return(nil, gin.H{"error": "ErrDatabaseFailure"})

		req, _ := http.NewRequest("POST", "/sign_up", strings.NewReader(`{"email":"test@example.com","name":"test","password":"Password123$"}`))
		// reqBody := `{"email": "test@example.com", "name": "test", "password": "Password123$"}`
		// req, _ := http.NewRequest("POST", "/sign_up", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/sign_up", userHandler.SignUpHandler())

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to create user", resp["error"])
	})

	t.Run("error: JWT generation fails", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		mockJWTGen := &mocks.MockJwtGenerator{}
		userHandler := user.NewHandler(mockClient, mockJWTGen)
		Email := "test@example.com"
		mockUser := &ent.User{ID: 1, Name: "test", Email: &Email}

		mockClient.On("Create", mock.Anything, "test@example.com", "test", mock.Anything).
			Return(mockUser, nil)
		mockJWTGen.On("GenerateJWT", "1").Return("", errors.New("JWT generation error"))

		req, _ := http.NewRequest("POST", "/sign_up", strings.NewReader(`{"email":"test@example.com","name":"test","password":"Password123$"}`))

		// reqBody := `{"email": "test@example.com", "name": "test", "password": "Password123$"}`
		// req, _ := http.NewRequest("POST", "/sign_up", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/sign_up", userHandler.SignUpHandler())

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to generate token", resp["error"])
	})

	t.Run("success: user created and token generated", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		mockJWTGen := &mocks.MockJwtGenerator{}
		userHandler := user.NewHandler(mockClient, mockJWTGen)
		// 正常なリクエストデータ
		reqData := models.SignUpRequest{
			Email:    "test@example.com",
			Name:     "TestUser",
			Password: "Secure123!",
		}
		reqBody, _ := json.Marshal(reqData)

		mockClient.On("Create", mock.Anything, reqData.Email, reqData.Name, mock.Anything).
			Return(&ent.User{ID: 1, Name: reqData.Name, Email: &reqData.Email}, nil)
		mockJWTGen.On("GenerateJWT", "1").Return("mocked_jwt_token", nil)

		req, _ := http.NewRequest(http.MethodPost, "/sign_up", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		userHandler.SignUpHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		responseData := map[string]string{}
		err := json.Unmarshal(w.Body.Bytes(), &responseData)
		assert.NoError(t, err)
		assert.Equal(t, "mocked_jwt_token", responseData["token"])

		mockClient.AssertExpectations(t)
		mockJWTGen.AssertExpectations(t)
	})

	t.Run("TestSignUpHandler_InvalidRequest", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockClient := new(user_mocks.MockUsecase)
		mockJWTGen := &mocks.MockJwtGenerator{}

		handler := user.NewHandler(mockClient, mockJWTGen)

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
	})
}
