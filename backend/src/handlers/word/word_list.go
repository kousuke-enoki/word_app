package word

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"word_app/backend/src/models"
	"word_app/backend/src/validators/word"

	"github.com/gin-gonic/gin"
)

func (h *WordHandler) WordListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		req, err := h.parseWordListRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// バリデーション
		validationErrors := word.ValidateWordListRequest(req)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}

		// サービスの呼び出し
		if req.SortBy == "register" {
			response, err := h.wordService.GetRegisteredWords(ctx, req)
			if err == nil {
				c.JSON(http.StatusOK, response)
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			response, err := h.wordService.GetWords(ctx, req)
			if err == nil {
				c.JSON(http.StatusOK, response)
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}
}

func (h *WordHandler) parseWordListRequest(c *gin.Context) (*models.WordListRequest, error) {
	// クエリパラメータの取得
	search := c.Query("search")
	sortBy := c.DefaultQuery("sortBy", "id")
	order := c.DefaultQuery("order", "asc")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		return nil, errors.New("Invalid 'page' query parameter: must be a positive integer")
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		return nil, errors.New("Invalid 'limit' query parameter: must be a positive integer")
	}

	// ユーザーIDをコンテキストから取得
	userID, exists := c.Get("userID")
	if !exists {
		return nil, errors.New("userID not found in context")
	}

	// userIDの型チェック
	userIDInt, ok := userID.(int)
	if !ok {
		return nil, errors.New("invalid userID type")
	}

	// リクエストオブジェクトを構築
	req := &models.WordListRequest{
		UserID: userIDInt,
		Search: search,
		SortBy: sortBy,
		Order:  order,
		Page:   page,
		Limit:  limit,
	}

	return req, nil
}
