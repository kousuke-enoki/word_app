package word

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *WordHandler) CreateWordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// リクエストを解析
		req, err := h.parseWordNewRequest(c)
		if err != nil {
			logrus.Errorf("Failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// サービス層にリクエストを渡して処理
		response, err := h.wordService.CreateWord(ctx, req)
		if err != nil {
			logrus.Errorf("Failed to create word: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create word"})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// リクエスト構造体を解析
func (h *WordHandler) parseWordNewRequest(c *gin.Context) (*models.CreateWordRequest, error) {
	var req models.CreateWordRequest

	// JSONリクエストをバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		return nil, err
	}

	// バリデーション
	if err := validateCreateWordRequest(&req); err != nil {
		logrus.Errorf("Validation error: %v", err)
		return nil, err
	}

	logrus.Infof("Parsed request: %+v", req)

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

	// コンテキストから取得したuserIDをリクエストに設定
	req.UserID = userIDInt
	logrus.Infof("Final parsed request with userID: %+v", req)

	return &req, nil
}

// バリデーション関数
func validateCreateWordRequest(req *models.CreateWordRequest) error {
	// 半角アルファベットのみの正規表現
	wordNameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)

	// 日本語（ひらがな、カタカナ、漢字）と記号「~」のみの正規表現
	japaneseMeanRegex := regexp.MustCompile(`^[ぁ-んァ-ヶ一-龠々ー～]+$`)

	// word.nameの検証
	if !wordNameRegex.MatchString(req.Name) {
		return errors.New("word.name must contain only alphabetic characters")
	}

	// WordInfosの検証
	if len(req.WordInfos) < 1 || len(req.WordInfos) > 10 {
		return errors.New("wordInfos must contain between 1 and 10 items")
	}

	// PartOfSpeechIDの重複チェック用マップ
	partOfSpeechIDMap := make(map[int]bool)

	for _, wordInfo := range req.WordInfos {
		// JapaneseMeansの検証
		if len(wordInfo.JapaneseMeans) < 1 || len(wordInfo.JapaneseMeans) > 10 {
			return errors.New("japaneseMeans must contain between 1 and 10 items")
		}

		for _, mean := range wordInfo.JapaneseMeans {
			if !japaneseMeanRegex.MatchString(mean.Name) {
				return errors.New("japaneseMean.name must contain only Japanese characters and the '~' symbol")
			}
		}

		// PartOfSpeechIDの重複検証
		if partOfSpeechIDMap[wordInfo.PartOfSpeechID] {
			return errors.New("duplicate PartOfSpeechID found in WordInfos")
		}
		partOfSpeechIDMap[wordInfo.PartOfSpeechID] = true
	}

	return nil
}
