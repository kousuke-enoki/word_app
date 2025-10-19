package word

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) ShowHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		req, err := h.parseWordShowRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// サービス層からデータを取得
		response, err := h.wordService.GetWordDetails(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h *Handler) parseWordShowRequest(c *gin.Context) (*models.WordShowRequest, error) {
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
	req := &models.WordShowRequest{
		WordID: wordID,
		UserID: userID,
	}

	logrus.Infof("Final parsed request: %+v", req)
	return req, nil
}
