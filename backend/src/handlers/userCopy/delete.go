// src/handlers/user/delete_handler.go
package user

import (
	"context"
	"net/http"
	"strconv"

	user_service "word_app/backend/src/service/user"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// 操作対象者
		// editorID はミドルウェアで設定済み前提
		editorIDAny, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		editorID, _ := editorIDAny.(int)

		// 削除対象IDの取得
		targetID, err := strconv.Atoi(c.Param("id"))
		if err != nil || targetID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		// サービス呼び出し
		err = h.userClient.Delete(ctx, editorID, targetID)
		if err != nil {
			switch err {
			case user_service.ErrUnauthorized:
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			case user_service.ErrUserNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			case user_service.ErrDatabaseFailure:
				fallthrough
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}
