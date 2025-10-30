package word

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) RegisterHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()
		req, err := h.parseRequest(c, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// サービス層からデータを取得
		response, err := h.wordService.RegisterWords(ctx, req)
		if err != nil {
			httperr.Write(c, err) // apperrorをそのまま返す
			return
		}

		c.JSON(http.StatusOK, response)
	})
}

func (h *Handler) parseRequest(c *gin.Context, userID int) (*models.RegisterWordRequest, error) {
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

	// JSONを`RegisterWordRequest`構造体にバインド
	var req models.RegisterWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("invalid JSON format: " + err.Error())
	}

	// コンテキストから取得したuserIDをリクエストに設定
	req.UserID = userID
	logrus.Infof("Final parsed request with userID: %+v", req)

	return &req, nil
}
