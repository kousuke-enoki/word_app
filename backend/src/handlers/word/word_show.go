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

func (h *WordHandler) WordShowHandler() gin.HandlerFunc {
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

func (h *WordHandler) parseWordShowRequest(c *gin.Context) (*models.WordShowRequest, error) {
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
	req := &models.WordShowRequest{
		WordID: wordID,
		UserID: userIDInt,
	}

	logrus.Infof("Final parsed request: %+v", req)
	return req, nil
}
