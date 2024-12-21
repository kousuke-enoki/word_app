package word

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"word_app/backend/src/models"
	"word_app/backend/src/validators/word"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *WordHandler) SaveMemoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		req, err := h.parseSaveMemoRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErrors := word.ValidateSaveMemo(req)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}

		// サービス層からデータを取得
		response, err := h.wordService.SaveMemo(ctx, req.WordID, req.UserID, req.Memo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h *WordHandler) parseSaveMemoRequest(c *gin.Context) (*models.SaveMemoRequest, error) {
	// リクエストボディが空の場合をチェック
	if c.Request.Body == nil {
		return nil, errors.New("request body is missing")
	}

	// リクエストボディを一旦読み込みつつ、後で再利用可能にする
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, errors.New("failed to read request body")
	}

	// GinのShouldBindJSONを正しく動作させるため、リクエストボディを再設定
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// JSONを`SaveMemoRequest`構造体にバインド
	var req models.SaveMemoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("invalid JSON format: " + err.Error())
	}

	// 必要に応じて追加処理（例: ユーザーIDをコンテキストから取得）
	userID, exists := c.Get("userID")
	if !exists {
		return nil, errors.New("unauthorized: userID not found in context")
	}

	// userIDの型チェック
	userIDInt, ok := userID.(int)
	if !ok {
		return nil, errors.New("invalid userID type")
	}

	// コンテキストから取得したuserIDをリクエストに設定
	req.UserID = userIDInt
	logrus.Infof("Final parsed request with userID: %+v", req)

	return &models.SaveMemoRequest{
		WordID: req.WordID,
		UserID: req.UserID,
		Memo:   req.Memo,
	}, nil
}
