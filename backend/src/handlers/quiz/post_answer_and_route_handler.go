package quiz

import (
	"context"
	"errors"
	"net/http"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *QuizHandler) PostAnswerAndRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// リクエストを解析
		answerReq, err := h.parsePostAnswerQuizRequest(c)
		if err != nil {
			logrus.Errorf("Failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ユーザーIDをコンテキストから取得
		userID, exists := c.Get("userID")
		if !exists {
			logrus.Error(errors.New("unauthorized: userID not found in context"))
			return
		}

		// userIDの型チェック
		userIDInt, ok := userID.(int)
		if !ok {
			logrus.Error(errors.New("invalid userID type"))
			return
		}

		// サービス層にリクエストを渡して処理
		response, err := h.quizService.SubmitAnswerAndRoute(ctx, userIDInt, answerReq)
		if err != nil {
			logrus.Errorf("Failed to create quiz: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz"})
			return
		}

		// // リクエストを解析
		// getReq, err := h.parseGetQuizRequest(c)
		// if err != nil {
		// 	logrus.Errorf("Failed to parse request: %v", err)
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }

		// getQuizResponse, err := h.quizService.GetQuiz(ctx, userIDInt, getReq)
		// if err != nil {
		// 	logrus.Errorf("Failed to create quiz: %v", err)
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz"})
		// 	return
		// }
		logrus.Info(response)

		c.JSON(http.StatusOK, response)
	}
}

// リクエスト構造体を解析
func (h *QuizHandler) parsePostAnswerQuizRequest(c *gin.Context) (*models.PostAnswerQuestionRequest, error) {
	var req models.PostAnswerQuestionRequest

	// JSONリクエストをバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		return nil, err
	}

	return &req, nil
}

// リクエスト構造体を解析
func (h *QuizHandler) parseGetQuizRequest(c *gin.Context) (*models.GetQuizRequest, error) {
	var req models.GetQuizRequest

	// JSONリクエストをバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		return nil, err
	}

	return &req, nil
}
