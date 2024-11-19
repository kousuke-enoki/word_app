package word

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *WordHandler) WordShowHandler(c *gin.Context) {
	ctx := context.Background()

	// パラメータの取得と検証
	wordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	// サービス層からデータを取得
	response, err := h.wordService.GetWordDetails(ctx, wordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
