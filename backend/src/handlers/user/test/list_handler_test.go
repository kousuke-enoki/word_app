package user_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"word_app/backend/src/handlers/user"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/mocks"
	user_mocks "word_app/backend/src/mocks/usecase/user"
	"word_app/backend/src/models"
	user_usecase "word_app/backend/src/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
  テストで網羅する分岐
  1) 200: 正常（デフォルト値適用）
  2) 200: 正常（任意クエリ: search/sortBy/order/page/limit をフォワード）
  3) 500: サービス層エラー（GetUsers が error を返す）
  4) 400: userID 不在
  5) 400: userID 型不正（string を入れる）
  6) 400: page 不正（0 / 非数）
  7) 400: limit 不正（0 / 非数）
*/

func newRouterWithUserID(h *user.UserHandler, userID any) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// userID をコンテキストに詰めるミドルウェア
	if userID != nil {
		r.Use(func(c *gin.Context) {
			var uid int
			switch v := userID.(type) {
			case int:
				uid = v
			default:
				// 型変換できない場合は何もしない
				c.Next()
				return
			}
			p := models.Principal{
				UserID:  uid,
				IsAdmin: false,
				IsRoot:  false,
				IsTest:  false,
			}
			jwt.SetPrincipal(c, p)
			c.Next()
		})
	}
	r.GET("/users", h.ListHandler())
	return r
}

func performGet(r *gin.Engine, rawURL string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, rawURL, nil)
	r.ServeHTTP(w, req)
	return w
}

func TestUserListHandler_AllPaths(t *testing.T) {
	t.Run("200 OK - defaults applied", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		mockJWTGen := &mocks.MockJwtGenerator{}
		h := user.NewHandler(mockClient, mockJWTGen)

		// 期待されるリクエスト値（デフォルト）
		expected := &user_usecase.ListUsersInput{
			ViewerID: 1,
			Search:   "",
			SortBy:   "name",
			Order:    "asc",
			Page:     1,
			Limit:    10,
		}

		// 引数一致：UserListRequest の全フィールドを確認
		argMatcher := mock.MatchedBy(func(req user_usecase.ListUsersInput) bool {
			return req.ViewerID == expected.ViewerID &&
				req.Search == expected.Search &&
				req.SortBy == expected.SortBy &&
				req.Order == expected.Order &&
				req.Page == expected.Page &&
				req.Limit == expected.Limit
		})

		Email := "alice@example.com"
		mockResp := &user_usecase.UserListResponse{
			Users: []models.User{
				{ID: 10, Name: "Alice", Email: &Email, IsAdmin: true, IsSettedPassword: true, IsLine: true},
			},
			TotalPages: 3,
		}
		mockClient.
			On("ListUsers", mock.Anything, argMatcher).
			Return(mockResp, nil)

		r := newRouterWithUserID(h, 1)
		w := performGet(r, "/users") // クエリ無し → デフォルト適用

		assert.Equal(t, http.StatusOK, w.Code)
		var got user_usecase.UserListResponse
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, mockResp.TotalPages, got.TotalPages)
		assert.Len(t, got.Users, 1)
		assert.Equal(t, 10, got.Users[0].ID)
		mockClient.AssertExpectations(t)
	})

	t.Run("200 OK - forwards query params", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		h := user.NewHandler(mockClient, nil)

		q := url.Values{}
		q.Set("search", "bob")
		q.Set("sortBy", "role")
		q.Set("order", "desc")
		q.Set("page", "2")
		q.Set("limit", "30")

		expected := &user_usecase.ListUsersInput{
			ViewerID: 1,
			Search:   "bob",
			SortBy:   "role",
			Order:    "desc",
			Page:     2,
			Limit:    30,
		}

		argMatcher := mock.MatchedBy(func(req user_usecase.ListUsersInput) bool {
			return req.ViewerID == expected.ViewerID &&
				req.Search == expected.Search &&
				req.SortBy == expected.SortBy &&
				req.Order == expected.Order &&
				req.Page == expected.Page &&
				req.Limit == expected.Limit
		})

		Email := "bob@example.com"
		mockResp := &user_usecase.UserListResponse{
			Users: []models.User{
				{ID: 20, Name: "Bob", Email: &Email, IsAdmin: false, IsRoot: true, IsSettedPassword: true, IsLine: false},
			},
			TotalPages: 7,
		}
		mockClient.
			On("ListUsers", mock.Anything, argMatcher).
			Return(mockResp, nil)

		r := newRouterWithUserID(h, 1)
		w := performGet(r, "/users?"+q.Encode())

		assert.Equal(t, http.StatusOK, w.Code)
		var got user_usecase.UserListResponse
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, 7, got.TotalPages)
		assert.Len(t, got.Users, 1)
		assert.Equal(t, "Bob", got.Users[0].Name)
		mockClient.AssertExpectations(t)
	})

	t.Run("500 - service error from ListUsers", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		h := user.NewHandler(mockClient, nil)

		anyReq := mock.AnythingOfType("user.ListUsersInput")
		mockClient.
			On("ListUsers", mock.Anything, anyReq).
			Return((*user_usecase.UserListResponse)(nil), errors.New("db down"))

		r := newRouterWithUserID(h, 1)
		w := performGet(r, "/users?search=x")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var got map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "internal error", got["error"])
		mockClient.AssertExpectations(t)
	})

	t.Run("400 - userID missing in context", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		h := user.NewHandler(mockClient, nil)

		// userID をコンテキストに入れない
		r := newRouterWithUserID(h, nil)
		w := performGet(r, "/users")

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var got map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "unauthorized", got["error"])

		// ListUsers は呼ばれない
		mockClient.AssertNotCalled(t, "ListUsers", mock.Anything, mock.Anything)
	})

	t.Run("400 - userID type invalid", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		h := user.NewHandler(mockClient, nil)

		// userID を string 型にして型不正を誘発
		r := newRouterWithUserID(h, "1")
		w := performGet(r, "/users")

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var got map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "unauthorized", got["error"])
		mockClient.AssertNotCalled(t, "ListUsers", mock.Anything, mock.Anything)
	})

	t.Run("400 - page invalid (zero)", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		h := user.NewHandler(mockClient, nil)

		r := newRouterWithUserID(h, 1)
		w := performGet(r, "/users?page=0")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var got map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "invalid 'page' query parameter: must be a positive integer", got["error"])
		mockClient.AssertNotCalled(t, "ListUsers", mock.Anything, mock.Anything)
	})

	t.Run("400 - page invalid (nan)", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		h := user.NewHandler(mockClient, nil)

		r := newRouterWithUserID(h, 1)
		w := performGet(r, "/users?page=abc")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var got map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "invalid 'page' query parameter: must be a positive integer", got["error"])
		mockClient.AssertNotCalled(t, "ListUsers", mock.Anything, mock.Anything)
	})

	t.Run("400 - limit invalid (zero)", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		h := user.NewHandler(mockClient, nil)

		r := newRouterWithUserID(h, 1)
		w := performGet(r, "/users?limit=0")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var got map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "invalid 'limit' query parameter: must be a positive integer", got["error"])
		mockClient.AssertNotCalled(t, "ListUsers", mock.Anything, mock.Anything)
	})

	t.Run("400 - limit invalid (nan)", func(t *testing.T) {
		mockClient := new(user_mocks.MockUsecase)
		h := user.NewHandler(mockClient, nil)

		r := newRouterWithUserID(h, 1)
		w := performGet(r, "/users?limit=abc")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var got map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "invalid 'limit' query parameter: must be a positive integer", got["error"])
		mockClient.AssertNotCalled(t, "ListUsers", mock.Anything, mock.Anything)
	})
}
