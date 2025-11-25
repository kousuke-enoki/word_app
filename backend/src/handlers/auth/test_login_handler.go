// src/handlers/auth/test_login.go
package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"net"
	"net/http"
	"strconv"
	"strings"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
)

const testLoginRoute = "/users/auth/test-login"

func HashUserAgent(ua string) string {
	h := sha1.Sum([]byte(ua))
	return hex.EncodeToString(h[:])
}

func (h *AuthHandler) TestLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		ip := clientIP(c)
		ua := c.Request.UserAgent()
		if ua == "" {
			ua = "unknown"
		}
		uaHash := HashUserAgent(ua)

		jump := c.DefaultQuery("jump", "quiz")

		out, lastPayload, retryAfter, err := h.AuthUsecase.TestLoginWithRateLimit(ctx, ip, uaHash, testLoginRoute, jump)
		if err != nil {
			// TooMany の場合だけ 429 + Retry-After
			if apperror.IsKind(err, apperror.TooManyRequests) {
				if retryAfter > 0 {
					c.Header("Retry-After", strconv.Itoa(retryAfter))
				}
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
				return
			}
			httperr.Write(c, err)
			return
		}
		if lastPayload != nil {
			// レート制限超過時のキャッシュレスポンスの場合、Retry-Afterヘッダーを設定
			if retryAfter > 0 {
				c.Header("Retry-After", strconv.Itoa(retryAfter))
				c.Header("X-Rate-Limit-Exceeded", "true") // レート制限超過フラグ
			}
			c.Data(http.StatusOK, "application/json", lastPayload)
			return
		}
		c.JSON(http.StatusOK, out)
	}
}

func clientIP(c *gin.Context) string {
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			if ip := strings.TrimSpace(parts[0]); ip != "" {
				return ip
			}
		}
	}
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return host
}
