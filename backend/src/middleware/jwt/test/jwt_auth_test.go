// jwt_auth_test.go
package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthenticator is a mock for jwt.Authenticator
type MockAuthenticator struct {
	mock.Mock
}

func (m *MockAuthenticator) Authenticate(ctx context.Context, raw string) (models.Principal, error) {
	args := m.Called(ctx, raw)
	if args.Get(0) == nil {
		return models.Principal{}, args.Error(1)
	}
	return args.Get(0).(models.Principal), args.Error(1)
}

func newMockAuthenticator(t *testing.T) *MockAuthenticator {
	return &MockAuthenticator{}
}

func newRouter(mw gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(mw)
	r.GET("/ping", func(c *gin.Context) {
		p, ok := jwt.GetPrincipal(c)
		if !ok {
			c.JSON(200, gin.H{"error": "no principal"})
			return
		}
		c.JSON(200, gin.H{
			"userID":  p.UserID,
			"isAdmin": p.IsAdmin,
			"isRoot":  p.IsRoot,
			"isTest":  p.IsTest,
		})
	})
	return r
}

func TestAuthenticateMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("header missing → 401", func(t *testing.T) {
		mockAuth := newMockAuthenticator(t)
		mw := jwt.NewMiddleware(mockAuth).AuthenticateMiddleware()
		r := newRouter(mw)

		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "authorization header required", resp["error"])

		mockAuth.AssertExpectations(t)
	})

	t.Run("invalid token → 401", func(t *testing.T) {
		mockAuth := newMockAuthenticator(t)
		mockAuth.On("Authenticate", mock.Anything, "badtoken").
			Return(models.Principal{}, errors.New("token_invalid parse_error")).
			Once()
		mw := jwt.NewMiddleware(mockAuth).AuthenticateMiddleware()
		r := newRouter(mw)

		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("Authorization", "Bearer badtoken")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "token_invalid parse_error", resp["error"])

		mockAuth.AssertExpectations(t)
	})

	t.Run("valid token → 200 & context set", func(t *testing.T) {
		mockAuth := newMockAuthenticator(t)
		p := models.Principal{UserID: 7, IsAdmin: true, IsRoot: false, IsTest: false}
		mockAuth.On("Authenticate", mock.Anything, "good").
			Return(p, nil).
			Once()

		mw := jwt.NewMiddleware(mockAuth).AuthenticateMiddleware()
		r := newRouter(mw)

		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("Authorization", "Bearer good")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, float64(7), resp["userID"])
		assert.Equal(t, true, resp["isAdmin"])
		assert.Equal(t, false, resp["isRoot"])
		assert.Equal(t, false, resp["isTest"])

		mockAuth.AssertExpectations(t)
	})

	t.Run("empty bearer token → 401", func(t *testing.T) {
		mockAuth := newMockAuthenticator(t)
		mw := jwt.NewMiddleware(mockAuth).AuthenticateMiddleware()
		r := newRouter(mw)

		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("Authorization", "Bearer ")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "authorization header required", resp["error"])

		mockAuth.AssertExpectations(t)
	})

	t.Run("different user roles", func(t *testing.T) {
		testCases := []struct {
			name      string
			principal models.Principal
		}{
			{
				name:      "admin user",
				principal: models.Principal{UserID: 1, IsAdmin: true, IsRoot: false, IsTest: false},
			},
			{
				name:      "root user",
				principal: models.Principal{UserID: 2, IsAdmin: false, IsRoot: true, IsTest: false},
			},
			{
				name:      "test user",
				principal: models.Principal{UserID: 3, IsAdmin: false, IsRoot: false, IsTest: true},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mockAuth := newMockAuthenticator(t)
				mockAuth.On("Authenticate", mock.Anything, "token").
					Return(tc.principal, nil).
					Once()

				mw := jwt.NewMiddleware(mockAuth).AuthenticateMiddleware()
				r := newRouter(mw)

				req := httptest.NewRequest(http.MethodGet, "/ping", nil)
				req.Header.Set("Authorization", "Bearer token")
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, float64(tc.principal.UserID), resp["userID"])
				assert.Equal(t, tc.principal.IsAdmin, resp["isAdmin"])
				assert.Equal(t, tc.principal.IsRoot, resp["isRoot"])
				assert.Equal(t, tc.principal.IsTest, resp["isTest"])

				mockAuth.AssertExpectations(t)
			})
		}
	})
}

func TestNewMiddleware(t *testing.T) {
	t.Run("success - create middleware", func(t *testing.T) {
		mockAuth := newMockAuthenticator(t)

		mw := jwt.NewMiddleware(mockAuth)

		assert.NotNil(t, mw)
		assert.Equal(t, mockAuth, mw.JwtUsecase)
	})
}
