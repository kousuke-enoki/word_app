// wrap_test.go
package jwt

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWithUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - handler called with user ID", func(t *testing.T) {
		router := gin.New()

		router.GET("/test", func(c *gin.Context) {
			// Principalを設定
			p := models.Principal{
				UserID:  123,
				IsAdmin: false,
				IsRoot:  false,
				IsTest:  false,
			}
			c.Set("principalKey", p)

			// WithUserでハンドラをラップ
			jwt.WithUser(func(c *gin.Context, userID int) {
				c.JSON(http.StatusOK, gin.H{"userID": userID})
			})(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		body := w.Body.Bytes()
		assert.NotEmpty(t, body)
	})

	t.Run("error - unauthorized - no principal", func(t *testing.T) {
		router := gin.New()

		router.GET("/test", jwt.WithUser(func(c *gin.Context, userID int) {
			c.JSON(http.StatusOK, gin.H{"userID": userID})
		}))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		body := w.Body.Bytes()
		assert.NotEmpty(t, body)
	})

	t.Run("success - principal with different user IDs", func(t *testing.T) {
		testCases := []int{1, 42, 999, 12345}

		for _, userID := range testCases {
			t.Run("user ID", func(t *testing.T) {
				router := gin.New()

				router.GET("/test", func(c *gin.Context) {
					p := models.Principal{
						UserID:  userID,
						IsAdmin: false,
						IsRoot:  false,
						IsTest:  false,
					}
					c.Set("principalKey", p)

					jwt.WithUser(func(c *gin.Context, uid int) {
						assert.Equal(t, userID, uid)
						c.JSON(http.StatusOK, gin.H{"userID": uid})
					})(c)
				})

				w := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
			})
		}
	})

	t.Run("success - handler can abort", func(t *testing.T) {
		router := gin.New()

		router.GET("/test", func(c *gin.Context) {
			p := models.Principal{
				UserID:  1,
				IsAdmin: false,
				IsRoot:  false,
				IsTest:  false,
			}
			c.Set("principalKey", p)

			jwt.WithUser(func(c *gin.Context, userID int) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			})(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("success - handler returns user details", func(t *testing.T) {
		router := gin.New()

		router.GET("/test", func(c *gin.Context) {
			p := models.Principal{
				UserID:  456,
				IsAdmin: true,
				IsRoot:  false,
				IsTest:  false,
			}
			c.Set("principalKey", p)

			jwt.WithUser(func(c *gin.Context, userID int) {
				p, ok := jwt.GetPrincipal(c)
				assert.True(t, ok)
				c.JSON(http.StatusOK, gin.H{
					"userID":  userID,
					"isAdmin": p.IsAdmin,
				})
			})(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
