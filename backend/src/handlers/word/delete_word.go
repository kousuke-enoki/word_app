package word

import (
	"errors"
	"net/http"
	"strconv"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Word削除用ハンドラー
func (h *Handler) DeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		principal, ok := jwt.GetPrincipal(c)
		if !ok || !principal.IsAdmin {
			httperr.Write(c, apperror.Unauthorizedf("unauthorized", nil))
			return
		}

		req, err := h.parseDeleteWordRequest(c, principal.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 削除処理を呼び出し
		response, err := h.wordService.DeleteWord(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 成功レスポンス
		c.JSON(http.StatusOK, response)
	}
}

func (h *Handler) parseDeleteWordRequest(c *gin.Context, userID int) (*models.DeleteWordRequest, error) {
	// パラメータの取得と検証
	wordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid word ID")
	}

	// リクエストオブジェクトを構築
	req := &models.DeleteWordRequest{
		WordID: wordID,
		UserID: userID,
	}

	logrus.Infof("Final parsed request: %+v", req)
	return req, nil
}
