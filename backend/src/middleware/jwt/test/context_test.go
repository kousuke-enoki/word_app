// context_test.go
package jwt

import (
	"context"
	"net/http"
	"testing"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetPrincipal(t *testing.T) {
	t.Run("success - set principal", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(nil)
		p := models.Principal{
			UserID:  1,
			IsAdmin: true,
			IsRoot:  false,
			IsTest:  false,
		}

		jwt.SetPrincipal(c, p)

		value, exists := c.Get("principalKey")
		assert.True(t, exists)
		principal, ok := value.(models.Principal)
		assert.True(t, ok)
		assert.Equal(t, p.UserID, principal.UserID)
		assert.Equal(t, p.IsAdmin, principal.IsAdmin)
		assert.Equal(t, p.IsRoot, principal.IsRoot)
		assert.Equal(t, p.IsTest, principal.IsTest)
	})
}

func TestGetPrincipal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - get principal", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		p := models.Principal{
			UserID:  42,
			IsAdmin: false,
			IsRoot:  true,
			IsTest:  false,
		}
		c.Set("principalKey", p)

		principal, ok := jwt.GetPrincipal(c)

		assert.True(t, ok)
		assert.Equal(t, 42, principal.UserID)
		assert.Equal(t, false, principal.IsAdmin)
		assert.Equal(t, true, principal.IsRoot)
		assert.Equal(t, false, principal.IsTest)
	})

	t.Run("not found - principal not set", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		principal, ok := jwt.GetPrincipal(c)

		assert.False(t, ok)
		assert.Equal(t, models.Principal{}, principal)
	})

	t.Run("not found - wrong type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Set("principalKey", "wrong type")

		principal, ok := jwt.GetPrincipal(c)

		assert.False(t, ok)
		assert.Equal(t, models.Principal{}, principal)
	})
}

func TestUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - get user ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		p := models.Principal{
			UserID:  123,
			IsAdmin: false,
			IsRoot:  false,
			IsTest:  false,
		}
		c.Set("principalKey", p)

		userID, ok := jwt.UserID(c)

		assert.True(t, ok)
		assert.Equal(t, 123, userID)
	})

	t.Run("not found - principal not set", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		userID, ok := jwt.UserID(c)

		assert.False(t, ok)
		assert.Equal(t, 0, userID)
	})
}

func TestRequireUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - require user ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		p := models.Principal{
			UserID:  456,
			IsAdmin: false,
			IsRoot:  false,
			IsTest:  false,
		}
		c.Set("principalKey", p)

		userID, err := jwt.RequireUserID(c)

		assert.NoError(t, err)
		assert.Equal(t, 456, userID)
	})

	t.Run("error - principal not set", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		userID, err := jwt.RequireUserID(c)

		assert.Error(t, err)
		assert.Equal(t, jwt.ErrUnauthorized, err)
		assert.Equal(t, 0, userID)
	})
}

func TestErrUnauthorized(t *testing.T) {
	t.Run("HTTPError implements error interface", func(t *testing.T) {
		err := jwt.ErrUnauthorized

		assert.Error(t, err)
		assert.Equal(t, "unauthorized", err.Error())
		assert.Equal(t, http.StatusUnauthorized, err.Code)
		assert.Equal(t, "unauthorized", err.Msg)
	})
}

func TestWithPrincipal(t *testing.T) {
	t.Run("success - add principal to context", func(t *testing.T) {
		ctx := context.Background()
		p := models.Principal{
			UserID:  789,
			IsAdmin: true,
			IsRoot:  false,
			IsTest:  false,
		}

		ctxWithPrincipal := jwt.WithPrincipal(ctx, p)

		retrieved, ok := jwt.GetPrincipalFromContext(ctxWithPrincipal)
		assert.True(t, ok)
		assert.Equal(t, 789, retrieved.UserID)
		assert.Equal(t, true, retrieved.IsAdmin)
		assert.Equal(t, false, retrieved.IsRoot)
		assert.Equal(t, false, retrieved.IsTest)
	})

	t.Run("success - add principal with test user", func(t *testing.T) {
		ctx := context.Background()
		p := models.Principal{
			UserID:  999,
			IsAdmin: false,
			IsRoot:  false,
			IsTest:  true,
		}

		ctxWithPrincipal := jwt.WithPrincipal(ctx, p)

		retrieved, ok := jwt.GetPrincipalFromContext(ctxWithPrincipal)
		assert.True(t, ok)
		assert.Equal(t, 999, retrieved.UserID)
		assert.Equal(t, true, retrieved.IsTest)
	})
}

func TestGetPrincipalFromContext(t *testing.T) {
	t.Run("success - get principal from context", func(t *testing.T) {
		ctx := context.Background()
		p := models.Principal{
			UserID:  111,
			IsAdmin: true,
			IsRoot:  true,
			IsTest:  false,
		}
		ctx = jwt.WithPrincipal(ctx, p)

		principal, ok := jwt.GetPrincipalFromContext(ctx)

		assert.True(t, ok)
		assert.Equal(t, 111, principal.UserID)
		assert.Equal(t, true, principal.IsAdmin)
		assert.Equal(t, true, principal.IsRoot)
		assert.Equal(t, false, principal.IsTest)
	})

	t.Run("not found - principal not in context", func(t *testing.T) {
		ctx := context.Background()

		principal, ok := jwt.GetPrincipalFromContext(ctx)

		assert.False(t, ok)
		assert.Equal(t, models.Principal{}, principal)
	})

	t.Run("not found - wrong type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), jwt.PrincipalKey{}, "wrong type")

		principal, ok := jwt.GetPrincipalFromContext(ctx)

		assert.False(t, ok)
		assert.Equal(t, models.Principal{}, principal)
	})
}

func TestIsTestUser(t *testing.T) {
	t.Run("success - is test user", func(t *testing.T) {
		ctx := context.Background()
		p := models.Principal{
			UserID:  1,
			IsAdmin: false,
			IsRoot:  false,
			IsTest:  true,
		}
		ctx = jwt.WithPrincipal(ctx, p)

		isTest := jwt.IsTestUser(ctx)

		assert.True(t, isTest)
	})

	t.Run("success - is not test user", func(t *testing.T) {
		ctx := context.Background()
		p := models.Principal{
			UserID:  1,
			IsAdmin: false,
			IsRoot:  false,
			IsTest:  false,
		}
		ctx = jwt.WithPrincipal(ctx, p)

		isTest := jwt.IsTestUser(ctx)

		assert.False(t, isTest)
	})

	t.Run("not found - principal not in context", func(t *testing.T) {
		ctx := context.Background()

		isTest := jwt.IsTestUser(ctx)

		assert.False(t, isTest)
	})
}

func TestWithPrincipal_ContextChaining(t *testing.T) {
	t.Run("success - chain multiple principals", func(t *testing.T) {
		ctx := context.Background()
		p1 := models.Principal{
			UserID:  1,
			IsAdmin: false,
			IsRoot:  false,
			IsTest:  false,
		}
		ctx1 := jwt.WithPrincipal(ctx, p1)

		p2 := models.Principal{
			UserID:  2,
			IsAdmin: true,
			IsRoot:  false,
			IsTest:  false,
		}
		ctx2 := jwt.WithPrincipal(ctx1, p2)

		// 最後に設定したものが優先される
		retrieved, ok := jwt.GetPrincipalFromContext(ctx2)
		assert.True(t, ok)
		assert.Equal(t, 2, retrieved.UserID)
		assert.Equal(t, true, retrieved.IsAdmin)
	})
}
