package word

import (
	"errors"
	"net/http"
	"strconv"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) ShowHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()

		req, err := h.parseWordShowRequest(c, userID)
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
	})
}

func (h *Handler) parseWordShowRequest(c *gin.Context, userID int) (*models.WordShowRequest, error) {
	// パラメータの取得と検証
	wordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid word ID")
	}

	// リクエストオブジェクトを構築
	req := &models.WordShowRequest{
		WordID: wordID,
		UserID: userID,
	}

	logrus.Infof("Final parsed request: %+v", req)
	return req, nil
}
