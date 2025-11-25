package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/config"
	"word_app/backend/src/handlers/auth"
	jwt_mock "word_app/backend/src/mocks/infrastructure/jwt"
	auth_mock "word_app/backend/src/mocks/usecase/auth"
	auth_usecase "word_app/backend/src/usecase/auth"
	"word_app/backend/src/usecase/shared/ucerr"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestLoginHandler_RateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("200 OK - within rate limit, no cache", func(t *testing.T) {
		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth.NewHandler(mockUC, mockJWTGen, config)

		// 新しいメソッド署名に合わせる
		mockUC.On("TestLoginWithRateLimit",
			mock.Anything,            // ctx
			"192.168.1.1",            // ip
			mock.Anything,            // uaHash (SHA1ハッシュ値)
			"/users/auth/test-login", // route
			"quiz",                   // jump
		).Return(
			&auth_usecase.TestLoginOutput{
				Token:    "test_jwt_token",
				UserID:   42,
				UserName: "テストユーザー@12345678",
				Jump:     "quiz",
			},
			nil, // lastPayload
			0,   // retryAfter
			nil, // error
		)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-login?jump=quiz", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.1")
		req.Header.Set("User-Agent", "test-agent")
		c.Request = req

		h.TestLoginHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "test_jwt_token", resp["token"])
		assert.Equal(t, float64(42), resp["user_id"]) // JSONはfloat64になる
		assert.Equal(t, "quiz", resp["jump"])

		mockUC.AssertExpectations(t)
	})

	t.Run("200 OK - return cached response", func(t *testing.T) {
		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth.NewHandler(mockUC, mockJWTGen, config)

		// キャッシュされたレスポンスを返す
		cachedPayload := []byte(`{"token":"cached-token","user_id":999,"user_name":"Cached User","jump":"list"}`)
		mockUC.On("TestLoginWithRateLimit",
			mock.Anything,
			"192.168.1.1",
			mock.Anything,
			"/users/auth/test-login",
			"quiz",
		).Return(
			nil,           // result (キャッシュ時はnil)
			cachedPayload, // lastPayload
			0,             // retryAfter
			nil,           // error
		)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-login?jump=quiz", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.1")
		req.Header.Set("User-Agent", "test-agent")
		c.Request = req

		h.TestLoginHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.JSONEq(t, string(cachedPayload), w.Body.String())

		mockUC.AssertExpectations(t)
	})

	t.Run("429 - rate limit exceeded (no cache)", func(t *testing.T) {
		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth.NewHandler(mockUC, mockJWTGen, config)

		// レート制限超過エラーを返す
		mockUC.On("TestLoginWithRateLimit",
			mock.Anything,
			"192.168.1.1",
			mock.Anything,
			"/users/auth/test-login",
			"quiz",
		).Return(
			nil, // result
			nil, // lastPayload
			45,  // retryAfter (45秒後に再試行)
			ucerr.TooManyRequests("rate limited"),
		)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-login?jump=quiz", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.1")
		req.Header.Set("User-Agent", "test-agent")
		c.Request = req

		h.TestLoginHandler()(c)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Equal(t, "45", w.Header().Get("Retry-After"))

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "too many requests", resp["error"])

		mockUC.AssertExpectations(t)
	})

	t.Run("403 - test mode disabled", func(t *testing.T) {
		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth.NewHandler(mockUC, mockJWTGen, config)

		mockUC.On("TestLoginWithRateLimit",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(
			nil,
			nil,
			0,
			ucerr.Forbidden("test user mode is disabled"),
		)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-login?jump=quiz", nil)
		c.Request = req

		h.TestLoginHandler()(c)

		assert.Equal(t, http.StatusForbidden, w.Code)

		mockUC.AssertExpectations(t)
	})

	t.Run("500 - internal server error", func(t *testing.T) {
		mockUC := new(auth_mock.MockUsecase)
		mockJWTGen := new(jwt_mock.MockJWTGenerator)
		config := &config.Config{}
		h := auth.NewHandler(mockUC, mockJWTGen, config)

		mockUC.On("TestLoginWithRateLimit",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(
			nil,
			nil,
			0,
			ucerr.Internal("database error", nil),
		)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/auth/test-login?jump=quiz", nil)
		c.Request = req

		h.TestLoginHandler()(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockUC.AssertExpectations(t)
	})
}
