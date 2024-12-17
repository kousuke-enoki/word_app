package word

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Word削除用ハンドラー
func (h *WordHandler) DeleteWordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		req, err := h.parseDeleteWordRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 削除処理を呼び出し
		response, err := h.wordService.DeleteWord(ctx, req.UserID, req.WordID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 成功レスポンス
		c.JSON(http.StatusOK, response)
	}
}

func (h *WordHandler) parseDeleteWordRequest(c *gin.Context) (*models.DeleteWordRequest, error) {
	// パラメータの取得と検証
	wordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("Invalid word ID")
	}

	// ユーザーIDをコンテキストから取得
	userID, exists := c.Get("userID")
	if !exists {
		return nil, errors.New("unauthorized: userID not found in context")
	}

	// userIDの型チェック
	userIDInt, ok := userID.(int)
	if !ok {
		return nil, errors.New("invalid userID type")
	}

	// リクエストオブジェクトを構築
	req := &models.DeleteWordRequest{
		WordID: wordID,
		UserID: userIDInt,
	}

	logrus.Infof("Final parsed request: %+v", req)
	return req, nil
}
