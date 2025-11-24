package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/config"
	auth_handler "word_app/backend/src/handlers/auth"
	jwt_mock "word_app/backend/src/mocks/infrastructure/jwt"
	auth_mock "word_app/backend/src/mocks/usecase/auth"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestLogoutHandler_AllPaths(t *testing.T) {
	t.Run("204 OK - success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth_handler.NewHandler(mockUC, mockJWTGen, config)

		mockUC.On("TestLogout", mock.Anything, 1).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-logout", nil)
		c.Request = req

		// Principalをセット
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: true}
		c.Set("principalKey", p)

		h.TestLogoutHandler()(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
		mockUC.AssertExpectations(t)
	})

	t.Run("204 OK - success with different user ID", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth_handler.NewHandler(mockUC, mockJWTGen, config)

		mockUC.On("TestLogout", mock.Anything, 999).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-logout", nil)
		c.Request = req

		// Principalをセット
		p := models.Principal{UserID: 999, IsAdmin: false, IsRoot: false, IsTest: true}
		c.Set("principalKey", p)

		h.TestLogoutHandler()(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
		mockUC.AssertExpectations(t)
	})

	t.Run("403 - forbidden (not test user)", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth_handler.NewHandler(mockUC, mockJWTGen, config)

		mockUC.On("TestLogout", mock.Anything, 1).
			Return(apperror.Forbiddenf("only test user can be deleted via test-logout", nil))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-logout", nil)
		c.Request = req

		// Principalをセット
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		h.TestLogoutHandler()(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "only test user can be deleted via test-logout")
		mockUC.AssertExpectations(t)
	})

	t.Run("500 - internal server error (database error)", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth_handler.NewHandler(mockUC, mockJWTGen, config)

		mockUC.On("TestLogout", mock.Anything, 1).
			Return(apperror.Internalf("database error", nil))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-logout", nil)
		c.Request = req

		// Principalをセット
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: true}
		c.Set("principalKey", p)

		h.TestLogoutHandler()(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
		mockUC.AssertExpectations(t)
	})

	t.Run("404 - not found (user already deleted)", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth_handler.NewHandler(mockUC, mockJWTGen, config)

		mockUC.On("TestLogout", mock.Anything, 1).
			Return(apperror.NotFoundf("user not found", nil))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-logout", nil)
		c.Request = req

		// Principalをセット
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: true}
		c.Set("principalKey", p)

		h.TestLogoutHandler()(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "user not found")
		mockUC.AssertExpectations(t)
	})

	t.Run("401 - missing principal", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth_handler.NewHandler(mockUC, mockJWTGen, config)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-logout", nil)
		c.Request = req

		// Principalをセットしない

		h.TestLogoutHandler()(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
		mockUC.AssertNotCalled(t, "TestLogout", mock.Anything, mock.Anything)
	})

	t.Run("204 OK - success with admin test user", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth_handler.NewHandler(mockUC, mockJWTGen, config)

		mockUC.On("TestLogout", mock.Anything, 5).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-logout", nil)
		c.Request = req

		// Principalをセット (admin & test)
		p := models.Principal{UserID: 5, IsAdmin: true, IsRoot: false, IsTest: true}
		c.Set("principalKey", p)

		h.TestLogoutHandler()(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
		mockUC.AssertExpectations(t)
	})
}
