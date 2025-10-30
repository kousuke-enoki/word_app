package word

import (
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"
	"word_app/backend/src/validators/word"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		principal, ok := jwt.GetPrincipal(c)
		if !ok || !principal.IsAdmin {
			httperr.Write(c, apperror.Unauthorizedf("unauthorized", nil))
			return
		}

		// リクエストを解析
		req, err := h.parseCreateWordRequest(c, principal.UserID)
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
func (h *Handler) parseCreateWordRequest(c *gin.Context, userID int) (*models.CreateWordRequest, error) {
	var req models.CreateWordRequest

	// JSONリクエストをバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		return nil, err
	}

	// コンテキストから取得したuserIDをリクエストに設定
	req.UserID = userID
	logrus.Infof("Final parsed request with userID: %+v", req)

	return &req, nil
}
