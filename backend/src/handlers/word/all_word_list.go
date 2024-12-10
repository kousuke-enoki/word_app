package word

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *WordHandler) AllWordListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		req, err := h.parseAllWordListRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// サービスの呼び出し
		response, err := h.wordService.GetWords(ctx, req.UserID, req.Search, req.SortBy, req.Order, req.Page, req.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h *WordHandler) parseAllWordListRequest(c *gin.Context) (*models.AllWordListRequest, error) {
	// クエリパラメータの取得
	search := c.Query("search")
	sortBy := c.DefaultQuery("sortBy", "id")
	order := c.DefaultQuery("order", "asc")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		return nil, errors.New("Invalid query parameters")
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		return nil, errors.New("Invalid query parameters")
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
	req := &models.AllWordListRequest{
		UserID: userIDInt,
		Search: search,
		SortBy: sortBy,
		Order:  order,
		Page:   page,
		Limit:  limit,
	}

	logrus.Infof("Final parsed request: %+v", req)
	return req, nil
}
