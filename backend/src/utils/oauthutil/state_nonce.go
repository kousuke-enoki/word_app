package oauthutil

import (
	"github.com/gin-gonic/gin"
)

// cookie 名は固定文字列で OK
const (
	stateCookie     = "line_oauth_state"
	nonceCookie     = "line_oauth_nonce"
	cookieMaxAgeSec = 300 // 5 分
)

// newState / newNonce ---------------

func NewState(c *gin.Context) (string, error) {
	s, err := NewRandomString()
	if err != nil {
		return "", err
	}
	writeCookie(c, stateCookie, s)
	return s, nil
}

func NewNonce(c *gin.Context) (string, error) {
	n, err := NewRandomString()
	if err != nil {
		return "", err
	}
	writeCookie(c, nonceCookie, n)
	return n, nil
}

// loadNonce (state も同様に必要なら) -----

func LoadNonce(c *gin.Context) string {
	val, _ := c.Cookie(nonceCookie)
	// 取り出したら即失効させるのがベター
	writeCookie(c, nonceCookie, "", -1)
	return val
}

// 共通 cookie 書き込みヘルパ
func writeCookie(c *gin.Context, name, val string, maxAge ...int) {
	age := cookieMaxAgeSec
	if len(maxAge) > 0 {
		age = maxAge[0]
	}
	c.SetCookie(
		name, val, // name, value
		age, "/", "", // path, domain
		true, true, // secure, httpOnly
	)
}
