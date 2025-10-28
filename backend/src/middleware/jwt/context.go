// middleware/context.go
package jwt

import (
	"context"
	"net/http"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type Principal struct {
	UserID  int
	IsAdmin bool
	IsRoot  bool
	IsTest  bool
	// 将来の拡張: Scopes []string, TenantID, DeviceID など
}

func SetPrincipal(c *gin.Context, p models.Principal) {
	c.Set("principalKey", p)
}

func GetPrincipal(c *gin.Context) (models.Principal, bool) {
	v, ok := c.Get("principalKey")
	if !ok {
		return models.Principal{}, false
	}
	p, ok := v.(models.Principal)
	return p, ok
}

// `(int, bool)` を返す軽量版
func UserID(c *gin.Context) (int, bool) {
	p, ok := GetPrincipal(c)
	if !ok {
		return 0, false
	}
	return p.UserID, true
}

// 401 用のエラーを返す「Require」版
func RequireUserID(c *gin.Context) (int, error) {
	p, ok := GetPrincipal(c)
	if !ok {
		return 0, ErrUnauthorized // 好みで apperror.Unauthorizedf に
	}
	return p.UserID, nil
}

var ErrUnauthorized = &HTTPError{Code: http.StatusUnauthorized, Msg: "unauthorized"}

type HTTPError struct {
	Code int
	Msg  string
}

func (e *HTTPError) Error() string { return e.Msg }

type principalKey struct{}

// WithPrincipal adds Principal to standard context.Context.
// This allows service/usecase layers to access Principal without gin.Context.
func WithPrincipal(ctx context.Context, p models.Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, p)
}

// GetPrincipalFromContext extracts Principal from standard context.Context.
// Returns the Principal and true if found, otherwise returns zero value and false.
func GetPrincipalFromContext(ctx context.Context) (models.Principal, bool) {
	v := ctx.Value(principalKey{})
	if v == nil {
		return models.Principal{}, false
	}
	p, ok := v.(models.Principal)
	return p, ok
}

// IsTestUser checks if the current user is a test user based on context.
// Returns true if user is a test user, false otherwise.
// Returns false if Principal is not found in context (assume not test user for safety).
func IsTestUser(ctx context.Context) bool {
	p, ok := GetPrincipalFromContext(ctx)
	if !ok {
		return false
	}
	return p.IsTest
}

// // WithPrincipalFromGin extracts Principal from gin.Context and embeds it into standard context.Context.
// // This is a convenience function for handlers to pass Principal to service/usecase layers.
// func WithPrincipalFromGin(c *gin.Context, ctx context.Context) context.Context {
// 	p, ok := GetPrincipal(c)
// 	if !ok {
// 		return ctx
// 	}
// 	return WithPrincipal(ctx, p)
// }
