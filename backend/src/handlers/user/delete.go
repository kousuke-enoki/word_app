// src/handlers/user/delete_handler.go
package user

import (
	"net/http"
	"strconv"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/usecase/apperror"
	user_usecase "word_app/backend/src/usecase/user"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) DeleteHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, editorID int) {
		ctx := c.Request.Context()

		// 削除対象IDの取得
		targetID, err := strconv.Atoi(c.Param("id"))
		if err != nil || targetID <= 0 {
			httperr.Write(c, apperror.Validationf("invalid userID type", nil))
			return
		}
		in := user_usecase.DeleteUserInput{
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
	})
}
