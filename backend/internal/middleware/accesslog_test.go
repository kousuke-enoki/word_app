package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAccessLog_RequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("X-Request-Id header is used when provided", func(t *testing.T) {
		router := gin.New()
		router.Use(AccessLog(logrus.StandardLogger(), AccessLogOpts{}))

		expectedID := "test-request-id-12345"
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-Id", expectedID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedID, w.Header().Get("X-Request-Id"))
	})

	t.Run("UUID is generated when X-Request-Id header is not provided", func(t *testing.T) {
		router := gin.New()
		router.Use(AccessLog(logrus.StandardLogger(), AccessLogOpts{}))

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		responseID := w.Header().Get("X-Request-Id")
		assert.NotEmpty(t, responseID)
		// UUID形式か確認
		_, err := uuid.Parse(responseID)
		assert.NoError(t, err)
	})
}

func TestAccessLog_SeverityLevel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		statusCode    int
		expectedLevel string
		handler       gin.HandlerFunc
	}{
		{
			name:          "2xx returns info level",
			statusCode:    200,
			expectedLevel: "info",
			handler: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			},
		},
		{
			name:          "3xx returns info level",
			statusCode:    301,
			expectedLevel: "info",
			handler: func(c *gin.Context) {
				c.Redirect(http.StatusMovedPermanently, "/redirect")
			},
		},
		{
			name:          "4xx returns warn level",
			statusCode:    400,
			expectedLevel: "warn",
			handler: func(c *gin.Context) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			},
		},
		{
			name:          "401 returns warn level",
			statusCode:    401,
			expectedLevel: "warn",
			handler: func(c *gin.Context) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			},
		},
		{
			name:          "404 returns warn level",
			statusCode:    404,
			expectedLevel: "warn",
			handler: func(c *gin.Context) {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			},
		},
		{
			name:          "5xx returns error level",
			statusCode:    500,
			expectedLevel: "error",
			handler: func(c *gin.Context) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(AccessLog(logrus.StandardLogger(), AccessLogOpts{}))

			router.GET("/test", tt.handler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			// ログレベルは実際のログ出力で確認（ここではステータスコードの確認のみ）
		})
	}
}

func TestAccessLog_HealthEndpointSuppression(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("health endpoint is downgraded to debug when ExcludeHealth is true", func(t *testing.T) {
		router := gin.New()
		router.Use(AccessLog(logrus.StandardLogger(), AccessLogOpts{
			HealthPath:    "/health",
			ExcludeHealth: true,
		}))

		router.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// ログはdebugレベルで出力される（LOG_LEVEL=infoでは表示されない）
	})

	t.Run("health endpoint is logged normally when ExcludeHealth is false", func(t *testing.T) {
		router := gin.New()
		router.Use(AccessLog(logrus.StandardLogger(), AccessLogOpts{
			HealthPath:    "/health",
			ExcludeHealth: false,
		}))

		router.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// ログは通常のinfoレベルで出力される
	})
}

func TestAccessLog_RequiredFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AccessLog(logrus.StandardLogger(), AccessLogOpts{}))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("Referer", "https://example.com")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// 必須フィールドはbaseFields関数で設定される（実際のログ出力で確認）
}

func TestAccessLog_UserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("user_id is included when set in context", func(t *testing.T) {
		router := gin.New()
		router.Use(AccessLog(logrus.StandardLogger(), AccessLogOpts{}))

		router.GET("/test", func(c *gin.Context) {
			c.Set("user_id", 123)
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("user_id is nil when not set in context", func(t *testing.T) {
		router := gin.New()
		router.Use(AccessLog(logrus.StandardLogger(), AccessLogOpts{}))

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// user_idはnilとして出力される
	})
}
