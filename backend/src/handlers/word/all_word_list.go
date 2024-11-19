package word

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *WordHandler) AllWordListHandler(c *gin.Context) {
	ctx := context.Background()

	// クエリパラメータの取得
	search := c.Query("search")
	sortBy := c.DefaultQuery("sortBy", "id")
	order := c.DefaultQuery("order", "asc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// サービスの呼び出し
	words, totalCount, totalPages, err := h.wordService.GetWords(ctx, search, sortBy, order, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンスの作成
	c.JSON(http.StatusOK, gin.H{
		"words":      words,
		"totalPages": totalPages,
		"totalCount": totalCount,
	})
}
