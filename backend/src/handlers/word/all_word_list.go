package word

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"word_app/ent"
	"word_app/ent/word"

	"github.com/gin-gonic/gin"
)

// GetAllWords 単語を取得するための関数。検索、ソート、ページネーションに対応。
func AllWordListHandler(c *gin.Context, client *ent.Client) {
	ctx := context.Background()

	// クエリパラメータの取得
	search := c.Query("search")                             // 検索クエリ
	sortBy := c.DefaultQuery("sortBy", "name")              // ソート基準 (デフォルトは 'name')
	order := c.DefaultQuery("order", "asc")                 // ソート順 ('asc' か 'desc')
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // ページ番号
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // 1ページあたりの件数

	// クエリの作成
	query := client.Word.Query()

	// 検索機能
	if search != "" {
		query = query.Where(word.NameContains(search))
	}

	// ソート機能
	if order == "asc" {
		query = query.Order(ent.Asc(sortBy))
	} else {
		query = query.Order(ent.Desc(sortBy))
	}

	// ページネーション機能
	offset := (page - 1) * limit
	words, err := query.Offset(offset).Limit(limit).WithWordInfos().All(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch words"})
		return
	}
	log.Println("all_word_list", words)
	// レスポンスとして単語のリストを返す
	c.JSON(http.StatusOK, words)
}
