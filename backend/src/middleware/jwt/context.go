// middleware/context.go
package jwt

import (
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
