package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	auth_handler "word_app/backend/src/handlers/auth"
	jwt_mock "word_app/backend/src/mocks/infrastructure/jwt"
	auth_mock "word_app/backend/src/mocks/usecase/auth"
	"word_app/backend/src/usecase/apperror"
	auth_usecase "word_app/backend/src/usecase/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestLoginHandler_AllPaths(t *testing.T) {
	t.Run("200 OK - success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		mockUC.On("TestLogin", mock.Anything).
			Return(&auth_usecase.TestLoginOutput{
				UserID:   42,
				UserName: "テストユーザー@12345678",
				Jump:     "list",
			}, nil)
		mockJWTGen.On("GenerateJWT", "42").Return("test_jwt_token", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-login", nil)
		c.Request = req

		h.TestLoginHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "test_jwt_token", resp["token"])

		userData := resp["user"].(map[string]interface{})
		assert.Equal(t, float64(42), userData["id"])
		assert.Equal(t, "テストユーザー@12345678", userData["name"])
		assert.True(t, userData["isTest"].(bool))

		mockUC.AssertExpectations(t)
		mockJWTGen.AssertExpectations(t)
	})

	t.Run("403 - forbidden (test user mode disabled)", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		mockUC.On("TestLogin", mock.Anything).
			Return(nil, apperror.Forbiddenf("test user mode is disabled", nil))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-login", nil)
		c.Request = req

		h.TestLoginHandler()(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "test user mode is disabled")
		mockJWTGen.AssertNotCalled(t, "GenerateJWT", mock.Anything)
		mockUC.AssertExpectations(t)
	})

	t.Run("500 - internal server error (database error)", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		mockUC.On("TestLogin", mock.Anything).
			Return(nil, apperror.Internalf("database error", nil))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-login", nil)
		c.Request = req

		h.TestLoginHandler()(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
		mockJWTGen.AssertNotCalled(t, "GenerateJWT", mock.Anything)
		mockUC.AssertExpectations(t)
	})

	t.Run("400 - jwt generation failure", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		mockUC.On("TestLogin", mock.Anything).
			Return(&auth_usecase.TestLoginOutput{
				UserID:   42,
				UserName: "テストユーザー@12345678",
				Jump:     "list",
			}, nil)
		mockJWTGen.On("GenerateJWT", "42").Return("", assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-login", nil)
		c.Request = req

		h.TestLoginHandler()(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to generate token")
		mockUC.AssertExpectations(t)
		mockJWTGen.AssertExpectations(t)
	})

	t.Run("200 OK - success with different user IDs", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		mockUC.On("TestLogin", mock.Anything).
			Return(&auth_usecase.TestLoginOutput{
				UserID:   999,
				UserName: "テストユーザー@abcdefgh",
				Jump:     "bulk",
			}, nil)
		mockJWTGen.On("GenerateJWT", "999").Return("different_token", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-login", nil)
		c.Request = req

		h.TestLoginHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "different_token", resp["token"])

		userData := resp["user"].(map[string]interface{})
		assert.Equal(t, float64(999), userData["id"])
		assert.Equal(t, "テストユーザー@abcdefgh", userData["name"])

		mockUC.AssertExpectations(t)
		mockJWTGen.AssertExpectations(t)
	})
}
