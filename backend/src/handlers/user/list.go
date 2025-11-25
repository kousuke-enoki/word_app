package user

import (
	"net/http"
	"strconv"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/usecase/apperror"
	user_usecase "word_app/backend/src/usecase/user"
	"word_app/backend/src/validators/user"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) ListHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, viewerID int) {
		ctx := c.Request.Context()

		req, err := h.parseUserListRequest(c, viewerID)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		// バリデーション
		validationErrors := user.ValidateUserListRequest(req)
		if len(validationErrors) > 0 {
			httperr.Write(c, apperror.WithFieldErrors(apperror.Validation, "invalid input", validationErrors))
			return
		}

		// サービスの呼び出し
		var (
			resp *user_usecase.UserListResponse
		)
		resp, err = h.userUsecase.ListUsers(ctx, *req)
		if err != nil {
			httperr.Write(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	})
}

func (h *UserHandler) parseUserListRequest(c *gin.Context, viewerID int) (*user_usecase.ListUsersInput, error) {
	// クエリパラメータの取得
	search := c.Query("search")
	sortBy := c.DefaultQuery("sortBy", "name")
	order := c.DefaultQuery("order", "asc")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		return nil, apperror.Validationf("invalid 'page' query parameter: must be a positive integer", err)
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		return nil, apperror.Validationf("invalid 'limit' query parameter: must be a positive integer", err)
	}

	// リクエストオブジェクトを構築
	req := &user_usecase.ListUsersInput{
		ViewerID: viewerID,
		Search:   search,
		SortBy:   sortBy,
		Order:    order,
		Page:     page,
		Limit:    limit,
	}

	return req, nil
}
