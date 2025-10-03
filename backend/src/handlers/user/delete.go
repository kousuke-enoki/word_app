// src/handlers/user/delete_handler.go
package user

import (
	"context"
	"net/http"
	"strconv"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/interfaces/http/user"
	"word_app/backend/src/usecase/apperror"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		// 操作対象者
		// editorID はミドルウェアで設定済み前提
		editorID, err := contextutil.MustUserID(c)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		// 削除対象IDの取得
		targetID, err := strconv.Atoi(c.Param("id"))
		if err != nil || targetID <= 0 {
			httperr.Write(c, apperror.Validationf("invalid userID type", nil))
			return
		}
		in := user.DeleteUserInput{
			EditorID: editorID,
			TargetID: targetID,
		}
		// ユースケース呼び出し
		err = h.userUsecase.Delete(ctx, in)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}
