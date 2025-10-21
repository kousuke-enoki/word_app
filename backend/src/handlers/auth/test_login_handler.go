// src/handlers/auth/test_login.go
package auth

import (

	// "crypto/sha1"
	// "encoding/hex"
	"fmt"
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
)

type TestLoginQuery struct {
	Jump string `form:"jump"` // list|bulk|quiz
	Size int    `form:"size"` // e.g. 10
}

func (h *AuthHandler) TestLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// 1) query 受け取り（未指定はデフォルト）
		// var q TestLoginQuery
		// _ = c.ShouldBindQuery(&q)
		// if q.Jump == "" {
		// 	q.Jump = "list"
		// }
		// if q.Size <= 0 {
		// 	q.Size = 10
		// }

		// // 2) クライアント情報（IP/UA）をUsecaseに渡す（将来的なレート制限やlast_resultに活用）
		// ip := clientIP(c)
		// ua := c.Request.UserAgent()
		// uaHash := sha1hex(ua)

		// // 3) 入力DTO
		// in := auth_usecase.TestLoginInput{
		// 	Jump:      q.Jump,
		// 	Size:      q.Size,
		// 	IP:        ip,
		// 	UAHash:    uaHash,
		// 	Now:       time.Now(),
		// 	RequestID: uuid.NewString(),
		// }

		// 4) テストユーザー作成 & 付随処理
		out, err := h.AuthUsecase.TestLogin(ctx)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		// 5) JWT発行
		token, err := h.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", out.UserID))
		if err != nil {
			httperr.Write(c, apperror.Validationf("Failed to generate token", err))
			return
		}

		// 6) レスポンス
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":     out.UserID,
				"name":   out.UserName,
				"isTest": true,
			},
		})
	}
}

// func clientIP(c *gin.Context) string {
// 	// X-Forwarded-For → X-Real-IP → RemoteAddr の順で取得（APIGW/ALB想定）
// 	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
// 		return ip
// 	}
// 	if ip := c.GetHeader("X-Real-IP"); ip != "" {
// 		return ip
// 	}
// 	host, _, _ := netSplitHostPort(c.Request.RemoteAddr)
// 	return host
// }

// func netSplitHostPort(hostport string) (host, port string, err error) {
// 	// 依存増やしたくないので簡易版
// 	for i := len(hostport) - 1; i >= 0; i-- {
// 		if hostport[i] == ':' {
// 			return hostport[:i], hostport[i+1:], nil
// 		}
// 	}
// 	return hostport, "", nil
// }

// func sha1hex(s string) string {
// 	h := sha1.Sum([]byte(s))
// 	return hex.EncodeToString(h[:])
// }
