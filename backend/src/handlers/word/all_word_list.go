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

// AllWordListHandler 単語を取得するための関数。検索、ソート、ページネーションに対応。
func AllWordListHandler(c *gin.Context, client *ent.Client) {
	ctx := context.Background()
	log.Println("AllWordListHandler")

	// クエリパラメータの取得
	search := c.Query("search")                             // 検索クエリ
	sortBy := c.DefaultQuery("sortBy", "id")                // ソート基準 (デフォルトは 'id')
	order := c.DefaultQuery("order", "asc")                 // ソート順 ('asc' か 'desc')
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // ページ番号
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // 1ページあたりの件数

	// クエリの作成（Wordを基準に）
	query := client.Word.Query()

	// 検索機能 (Wordの名前で検索)
	if search != "" {
		query = query.Where(
			word.NameContains(search), // Wordのnameフィールドで検索
		)
	}
	log.Println(query)

	// ページネーション機能
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Wordに紐づくWordInfo, JapaneseMean, PartOfSpeechを取得
	query = query.
		WithWordInfos( // Wordに紐づくWordInfoを取得
			func(wiQuery *ent.WordInfoQuery) {
				wiQuery.WithJapaneseMeans().WithPartOfSpeech()
			},
		)

	// ソート機能
	if sortBy == "name" {
		// Word の name フィールドでソート
		if order == "asc" {
			query = query.Order(
				ent.Asc(word.FieldName), // Word の name で昇順ソート
			)
		} else {
			query = query.Order(
				ent.Desc(word.FieldName), // Word の name で降順ソート
			)
		}
	} else {
		// それ以外のフィールドでソート
		if order == "asc" {
			query = query.Order(ent.Asc(sortBy))
		} else {
			query = query.Order(ent.Desc(sortBy))
		}
	}
	log.Println(query)

	// クエリ実行
	words, err := query.All(ctx)
	if err != nil {
		log.Println("Error fetching words:", err) // エラーを詳細に出力
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch words"})
		return
	}

	// ログに取得したデータを表示
	log.Println("words_with_relations", words)

	// レスポンスとしてWordのリストを返す
	c.JSON(http.StatusOK, words)
}
