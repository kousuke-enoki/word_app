package word

import (
	"context"
	"net/http"
	"strconv"
	"word_app/ent"
	"word_app/ent/word"

	"github.com/gin-gonic/gin"
)

// WordResponse 構造体でレスポンスを定義
type WordResponse struct {
	Name               string     `json:"name"`
	WordInfos          []WordInfo `json:"wordInfos"`
	IsRegistered       bool       `json:"isRegistered"`
	TestCount          int        `json:"testCount"`
	CheckCount         int        `json:"checkCount"`
	RegistrationActive bool       `json:"registrationActive"`
	Memo               string     `json:"memo"`
}

type WordInfo struct {
	ID            int            `json:"id"`
	PartOfSpeech  PartOfSpeech   `json:"partOfSpeech"`
	JapaneseMeans []JapaneseMean `json:"japaneseMeans"`
}

type PartOfSpeech struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type JapaneseMean struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// WordShowHandler 単語を取得するための関数
func WordShowHandler(c *gin.Context, client *ent.Client) {
	ctx := context.Background()
	wordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	word, err := client.Word.
		Query().
		Where(word.ID(wordID)).
		WithWordInfos(func(wq *ent.WordInfoQuery) {
			wq.WithJapaneseMeans().WithPartOfSpeech()
		}).
		WithRegisteredWords().
		Only(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch word details"})
		return
	}

	// 登録済み情報を取得
	var isRegistered bool
	var testCount, checkCount int
	var memo string

	if len(word.Edges.RegisteredWords) > 0 {
		registeredWord := word.Edges.RegisteredWords[0]
		isRegistered = true
		testCount = registeredWord.TestCount
		checkCount = registeredWord.CheckCount
		if registeredWord.Memo != nil {
			memo = *registeredWord.Memo
		} else {
			memo = ""
		}
	}

	// WordInfosを変換
	wordInfos := make([]WordInfo, len(word.Edges.WordInfos))
	for i, wordInfo := range word.Edges.WordInfos {
		partOfSpeech := PartOfSpeech{
			ID:   wordInfo.Edges.PartOfSpeech.ID,
			Name: wordInfo.Edges.PartOfSpeech.Name,
		}
		japaneseMeans := make([]JapaneseMean, len(wordInfo.Edges.JapaneseMeans))
		for j, mean := range wordInfo.Edges.JapaneseMeans {
			japaneseMeans[j] = JapaneseMean{
				ID:   mean.ID,
				Name: mean.Name,
			}
		}
		wordInfos[i] = WordInfo{
			ID:            wordInfo.ID,
			PartOfSpeech:  partOfSpeech,
			JapaneseMeans: japaneseMeans,
		}
	}

	response := WordResponse{
		Name:               word.Name,
		WordInfos:          wordInfos,
		IsRegistered:       isRegistered,
		TestCount:          testCount,
		CheckCount:         checkCount,
		RegistrationActive: isRegistered,
		Memo:               memo,
	}

	c.JSON(http.StatusOK, response)
}
