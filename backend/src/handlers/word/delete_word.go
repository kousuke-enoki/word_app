package word

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Word削除用ハンドラー
func (h *WordHandler) DeleteWordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		// トークンから userID を取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// クエリパラメータから wordID を取得
		wordIDParam := c.Param("id")
		wordID, err := strconv.Atoi(wordIDParam)
		if err != nil || wordID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wordID"})
			return
		}

		// 削除処理を呼び出し
		response, err := h.wordService.DeleteWord(ctx, userID.(int), wordID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 成功レスポンス
		c.JSON(http.StatusOK, response)
	}
}
