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

func (h *WordHandler) WordNewHandler() gin.HandlerFunc {
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
func (h *WordHandler) parseWordNewRequest(c *gin.Context) (*models.WordCreateRequest, error) {
	var req models.WordCreateRequest

	// JSONリクエストをバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		return nil, err
	}

	// バリデーション
	if err := validateWordCreateRequest(&req); err != nil {
		logrus.Errorf("Validation error: %v", err)
		return nil, err
	}

	logrus.Infof("Parsed request: %+v", req)
	return &req, nil
}

// バリデーション関数
func validateWordCreateRequest(req *models.WordCreateRequest) error {
	// 半角アルファベットのみの正規表現
	wordNameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	// アルファベット以外の正規表現
	japaneseMeanRegex := regexp.MustCompile(`^[^\x00-\x7F]+$`)

	// word.nameの検証
	if !wordNameRegex.MatchString(req.Name) {
		return errors.New("word.name must contain only alphabetic characters")
	}

	// wordInfosの検証
	for _, wordInfo := range req.WordInfos {
		for _, mean := range wordInfo.JapaneseMeans {
			if !japaneseMeanRegex.MatchString(mean.Name) {
				return errors.New("japaneseMean.name must contain only non-alphabetic characters")
			}
		}
	}

	return nil
}
