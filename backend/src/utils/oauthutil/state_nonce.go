package oauthutil

import (
	"crypto/subtle"
	"net/http"

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

// 期待する state とクエリの state を比較（恒等時間比較）
// 検証の成否に関わらず Cookie は即失効させる
// const stateCookie = "line_oauth_state"

func VerifyState(c *gin.Context, got string) bool {
	want, err := c.Cookie(stateCookie)
	// 取得できなければ false
	if err != nil || got == "" {
		return false
	}
	// 一度使った state は無効化（削除）
	writeCookie(c, stateCookie, "", -1)
	// タイミング攻撃対策で比較
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
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
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    val,
		Path:     "/",
		MaxAge:   age,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		// Domain:  空(ホスト限定でOK。APIのホストに付く)
	})
}
