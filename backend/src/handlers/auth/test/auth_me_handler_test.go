package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	auth_handler "word_app/backend/src/handlers/auth"
	jwt_mock "word_app/backend/src/mocks/infrastructure/jwt"
	auth_mock "word_app/backend/src/mocks/usecase/auth"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthMeHandler_AllPaths(t *testing.T) {
	t.Run("200 OK - success with normal user", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		c.Request = req

		// Principalをセット
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		h.AuthMeHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.MeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.User.ID)
		assert.False(t, resp.User.IsAdmin)
		assert.False(t, resp.User.IsRoot)
		assert.False(t, resp.User.IsTest)
		mockUC.AssertNotCalled(t, "TestLogin", mock.Anything)
	})

	t.Run("200 OK - success with admin user", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		c.Request = req

		// Principalをセット
		p := models.Principal{UserID: 2, IsAdmin: true, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		h.AuthMeHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.MeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 2, resp.User.ID)
		assert.True(t, resp.User.IsAdmin)
		assert.False(t, resp.User.IsRoot)
		assert.False(t, resp.User.IsTest)
		mockUC.AssertNotCalled(t, "TestLogin", mock.Anything)
	})

	t.Run("200 OK - success with root user", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		c.Request = req

		// Principalをセット
		p := models.Principal{UserID: 3, IsAdmin: true, IsRoot: true, IsTest: false}
		c.Set("principalKey", p)

		h.AuthMeHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.MeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 3, resp.User.ID)
		assert.True(t, resp.User.IsAdmin)
		assert.True(t, resp.User.IsRoot)
		assert.False(t, resp.User.IsTest)
		mockUC.AssertNotCalled(t, "TestLogin", mock.Anything)
	})

	t.Run("200 OK - success with test user", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		c.Request = req

		// Principalをセット
		p := models.Principal{UserID: 4, IsAdmin: false, IsRoot: false, IsTest: true}
		c.Set("principalKey", p)

		h.AuthMeHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.MeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 4, resp.User.ID)
		assert.False(t, resp.User.IsAdmin)
		assert.False(t, resp.User.IsRoot)
		assert.True(t, resp.User.IsTest)
		mockUC.AssertNotCalled(t, "TestLogin", mock.Anything)
	})

	t.Run("401 - missing principal", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		c.Request = req

		// Principalをセットしない

		h.AuthMeHandler()(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
		mockUC.AssertNotCalled(t, "TestLogin", mock.Anything)
	})

	t.Run("401 - invalid principal type", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		h := auth_handler.NewHandler(mockUC, mockJWTGen)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		c.Request = req

		// 不正な型のPrincipalをセット
		c.Set("principalKey", "invalid")

		h.AuthMeHandler()(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
		mockUC.AssertNotCalled(t, "TestLogin", mock.Anything)
	})
}
