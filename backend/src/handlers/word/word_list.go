package word

import (
	"errors"
	"net/http"
	"strconv"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"
	"word_app/backend/src/validators/word"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()
		req, err := h.parseWordListRequest(c, userID)
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
		var (
			resp *models.WordListResponse
		)
		if req.SortBy == "register" {
			resp, err = h.wordService.GetRegisteredWords(ctx, req)
		} else {
			resp, err = h.wordService.GetWords(ctx, req)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
}

func (h *Handler) parseWordListRequest(c *gin.Context, userID int) (*models.WordListRequest, error) {
	// クエリパラメータの取得
	search := c.Query("search")
	sortBy := c.DefaultQuery("sortBy", "id")
	order := c.DefaultQuery("order", "asc")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		return nil, errors.New("invalid 'page' query parameter: must be a positive integer")
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		return nil, errors.New("invalid 'limit' query parameter: must be a positive integer")
	}

	// リクエストオブジェクトを構築
	req := &models.WordListRequest{
		UserID: userID,
		Search: search,
		SortBy: sortBy,
		Order:  order,
		Page:   page,
		Limit:  limit,
	}

	return req, nil
}
