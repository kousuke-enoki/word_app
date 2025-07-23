package word

import (
	"context"
	"errors"
	"net/http"

	"word_app/backend/src/models"
	"word_app/backend/src/utils/contextutil"
	"word_app/backend/src/validators/word"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		userRoles, err := contextutil.GetUserRoles(c)
		if err != nil || userRoles == nil || !userRoles.IsAdmin {
			if err == nil {
				err = errors.New("unauthorized: admin access required")
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// リクエストを解析
		req, err := h.parseCreateWordRequest(c)
		if err != nil {
			logrus.Errorf("Failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// バリデーション
		validationErrors := word.ValidateCreateWordRequest(req)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
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
func (h *Handler) parseCreateWordRequest(c *gin.Context) (*models.CreateWordRequest, error) {
	var req models.CreateWordRequest

	// JSONリクエストをバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		return nil, err
	}

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
