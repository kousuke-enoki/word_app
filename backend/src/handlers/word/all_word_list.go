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
		query = query.Where(word.NameContains(search))
	}

	// 総レコード数をカウント
	totalCount, err := query.Count(ctx)
	if err != nil {
		log.Println("Error counting words:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count words"})
		return
	}
	log.Println("totalCount")
	log.Println(totalCount)
	// ページネーション機能
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Wordに紐づくWordInfo, JapaneseMean, PartOfSpeechを取得
	query = query.WithWordInfos(
		func(wiQuery *ent.WordInfoQuery) {
			wiQuery.WithJapaneseMeans().WithPartOfSpeech()
		},
	)

	// ソート機能
	switch sortBy {
	case "name":
		if order == "asc" {
			query = query.Order(ent.Asc(word.FieldName))
		} else {
			query = query.Order(ent.Desc(word.FieldName))
		}
	// case "part_of_speech_name":
	// PartOfSpeechのnameフィールドでソート
	// query = query.Join("word_infos").
	// 	Join("part_of_speech").
	// 	OrderFunc(func(builder *sql.Selector) {
	// 		if order == "asc" {
	// 			builder.OrderBy(sql.Asc("part_of_speech.name"))
	// 		} else {
	// 			builder.OrderBy(sql.Desc("part_of_speech.name"))
	// 		}
	// 	})
	default:
		if order == "asc" {
			query = query.Order(ent.Asc(sortBy))
		} else {
			query = query.Order(ent.Desc(sortBy))
		}
	}

	// クエリ実行
	words, err := query.All(ctx)
	if err != nil {
		log.Println("Error fetching words:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch words"})
		return
	}

	// 総ページ数を計算
	totalPages := (totalCount + limit - 1) / limit
	log.Println("totalPages")
	log.Println(totalPages)
	// ログに取得したデータを表示
	log.Println("words_with_relations", words)

	// レスポンスとしてWordのリストと総ページ数を返す
	c.JSON(http.StatusOK, gin.H{
		"words":      words,
		"totalPages": totalPages,
		"totalCount": totalCount,
	})
}
