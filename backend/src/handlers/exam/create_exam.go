package exam

import (
	"context"
	"errors"
	"net/http"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *ExamHandler) CreateExamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// リクエストを解析
		req, err := h.parseCreateExamRequest(c)
		if err != nil {
			logrus.Errorf("Failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// // バリデーション
		// validationErrors := word.ValidateCreateExamRequest(req)
		// if len(validationErrors) > 0 {
		// 	c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		// 	return
		// }

		// サービス層にリクエストを渡して処理
		response, err := h.examService.CreateExam(ctx, req)
		if err != nil {
			logrus.Errorf("Failed to create exam: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exam"})
			return
		}
		logrus.Info(response)

		c.JSON(http.StatusOK, response)
	}
}

// リクエスト構造体を解析
func (h *ExamHandler) parseCreateExamRequest(c *gin.Context) (*models.CreateExamRequest, error) {
	var req models.CreateExamRequest

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
