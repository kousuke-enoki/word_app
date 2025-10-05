// src/handlers/user/test/detail_handler_test.go
package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	h "word_app/backend/src/handlers/user"
	user_mocks "word_app/backend/src/mocks/http/user"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/************ ヘルパー ************/

// ルータ生成（/me 用）。必要なら userID を注入するミドルウェアをルート登録前に設定
func newMeRouter(uc *user_mocks.MockUsecase, userID any) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if userID != nil {
		r.Use(func(c *gin.Context) {
			c.Set("userID", userID)
			c.Next()
		})
	}
	hd := h.NewHandler(uc, nil)
	r.GET("/me", hd.MeHandler())
	return r
}

// ルータ生成（/users/:id 用）
func newShowRouter(uc *user_mocks.MockUsecase, userID any) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if userID != nil {
		r.Use(func(c *gin.Context) {
			c.Set("userID", userID)
			c.Next()
		})
	}
	hd := h.NewHandler(uc, nil)
	r.GET("/users/:id", hd.ShowHandler())
	return r
}

func doGET(r *gin.Engine, url string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	r.ServeHTTP(w, req)
	return w
}

/************ MeHandler ************/

func TestMeHandler_AllPaths(t *testing.T) {
	t.Run("200 OK - success", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newMeRouter(uc, 101)

		uc.On("GetMyDetail", mock.Anything, 101).
			Return(&models.UserDetail{
				ID:      101,
				Name:    "Alice",
				Email:   nil,
				IsAdmin: true,
				IsRoot:  false,
				IsTest:  false,
			}, nil)

		w := doGET(r, "/me")

		assert.Equal(t, http.StatusOK, w.Code)
		var got models.UserDetail
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, 101, got.ID)
		assert.Equal(t, "Alice", got.Name)
		assert.True(t, got.IsAdmin)
		uc.AssertExpectations(t)
	})

	t.Run("401 - userID missing", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newMeRouter(uc, nil)

		w := doGET(r, "/me")

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"unauthorized: userID not found in context"}`, w.Body.String())
		uc.AssertNotCalled(t, "GetMyDetail", mock.Anything, mock.Anything)
	})

	t.Run("500 - internal error from usecase", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newMeRouter(uc, 7)

		uc.On("GetMyDetail", mock.Anything, 7).
			Return((*models.UserDetail)(nil), apperror.Internalf("internal error", nil))

		w := doGET(r, "/me")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"internal error"}`, w.Body.String())
		uc.AssertExpectations(t)
	})
}

/************ ShowHandler ************/

func TestShowHandler_AllPaths(t *testing.T) {
	t.Run("200 OK - success", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newShowRouter(uc, 9)

		// viewer=9, target=123
		uc.On("GetDetailByID", mock.Anything, 9, 123).
			Return(&models.UserDetail{
				ID:      123,
				Name:    "Bob",
				Email:   nil,
				IsAdmin: false,
				IsRoot:  true,
				IsTest:  false,
			}, nil)

		w := doGET(r, "/users/123")

		assert.Equal(t, http.StatusOK, w.Code)
		var got models.UserDetail
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, 123, got.ID)
		assert.Equal(t, "Bob", got.Name)
		assert.True(t, got.IsRoot)
		uc.AssertExpectations(t)
	})

	t.Run("401 - userID missing", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newShowRouter(uc, nil)

		w := doGET(r, "/users/1")

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"unauthorized: userID not found in context"}`, w.Body.String())
		uc.AssertNotCalled(t, "GetDetailByID", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("400 - invalid id (non-numeric)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newShowRouter(uc, 1)

		w := doGET(r, "/users/abc")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid id"}`, w.Body.String())
		uc.AssertNotCalled(t, "GetDetailByID", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("400 - invalid id (<= 0)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newShowRouter(uc, 1)

		w := doGET(r, "/users/0")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid id"}`, w.Body.String())
		uc.AssertNotCalled(t, "GetDetailByID", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("404 - not found (usecase returns NotFound)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newShowRouter(uc, 2)

		uc.On("GetDetailByID", mock.Anything, 2, 777).
			Return((*models.UserDetail)(nil), apperror.NotFoundf("user not found", nil))

		w := doGET(r, "/users/777")

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"user not found"}`, w.Body.String())
		uc.AssertExpectations(t)
	})

	t.Run("403 - forbidden (usecase returns Forbidden)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newShowRouter(uc, 3)

		uc.On("GetDetailByID", mock.Anything, 3, 4).
			Return((*models.UserDetail)(nil), apperror.Forbiddenf("forbidden", nil))

		w := doGET(r, "/users/4")

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.JSONEq(t, `{"error":"forbidden"}`, w.Body.String())
		uc.AssertExpectations(t)
	})

	t.Run("500 - internal error (usecase returns Internal)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newShowRouter(uc, 5)

		uc.On("GetDetailByID", mock.Anything, 5, 6).
			Return((*models.UserDetail)(nil), apperror.Internalf("internal error", nil))

		w := doGET(r, "/users/6")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"internal error"}`, w.Body.String())
		uc.AssertExpectations(t)
	})
}
