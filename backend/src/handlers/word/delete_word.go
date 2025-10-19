package word

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Word削除用ハンドラー
func (h *Handler) DeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		userRoles, err := contextutil.GetUserRoles(c)
		if err != nil || userRoles == nil || !userRoles.IsAdmin {
			if err == nil {
				err = errors.New("unauthorized: admin access required")
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		req, err := h.parseDeleteWordRequest(c)
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

func (h *Handler) parseDeleteWordRequest(c *gin.Context) (*models.DeleteWordRequest, error) {
	// パラメータの取得と検証
	wordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid word ID")
	}

	// ユーザーIDをコンテキストから取得
	userID, err := jwt.RequireUserID(c)
	if err != nil {
		return nil, errors.New("unauthorized: userID not found in context")
	}

	// リクエストオブジェクトを構築
	req := &models.DeleteWordRequest{
		WordID: wordID,
		UserID: userID,
	}

	logrus.Infof("Final parsed request: %+v", req)
	return req, nil
}
