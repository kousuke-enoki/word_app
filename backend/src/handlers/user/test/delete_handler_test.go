// src/handlers/user/test/delete_handler_test.go
package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	h "word_app/backend/src/handlers/user"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/mocks"
	"word_app/backend/src/models"
	user_mocks "word_app/backend/src/mocks/usecase/user"
	"word_app/backend/src/usecase/apperror"
	user_usecase "word_app/backend/src/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newDeleteRouterWithUserID(uc *user_mocks.MockUsecase, userID any) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// ここで Principal を入れるミドルウェアを"ルート登録前"に設定
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
	handler := h.NewHandler(uc, &mocks.MockJwtGenerator{})
	r.DELETE("/users/:id", handler.DeleteHandler())
	return r
}

func performDELETE(r *gin.Engine, url string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	r.ServeHTTP(w, req)
	return w
}

func TestDeleteHandler_AllPaths(t *testing.T) {
	t.Run("200 OK - success", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		// editorID=9 をミドルウェアで注入
		r := newDeleteRouterWithUserID(uc, 9)

		in := user_usecase.DeleteUserInput{EditorID: 9, TargetID: 123}
		uc.On("Delete", mock.Anything, in).Return(nil)

		w := performDELETE(r, "/users/123")

		assert.Equal(t, http.StatusOK, w.Code)
		var got map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "deleted", got["message"])
		uc.AssertExpectations(t)
	})

	t.Run("401 - userID missing", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		// userID を注入しない
		r := newDeleteRouterWithUserID(uc, nil)

		w := performDELETE(r, "/users/10")

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"unauthorized"}`, w.Body.String())
		uc.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("400 - invalid target id (non numeric)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		// userID は正常に入れて、ターゲットIDの検証まで到達させる
		r := newDeleteRouterWithUserID(uc, 1)

		w := performDELETE(r, "/users/abc")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid userID type"}`, w.Body.String())
		uc.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("400 - invalid target id (<= 0)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newDeleteRouterWithUserID(uc, 1)

		w := performDELETE(r, "/users/0")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid userID type"}`, w.Body.String())
		uc.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("404 - not found (usecase returns NotFound)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newDeleteRouterWithUserID(uc, 2)

		in := user_usecase.DeleteUserInput{EditorID: 2, TargetID: 555}
		uc.On("Delete", mock.Anything, in).Return(apperror.NotFoundf("user not found", nil))

		w := performDELETE(r, "/users/555")

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"user not found"}`, w.Body.String())
		uc.AssertExpectations(t)
	})

	t.Run("403 - forbidden (usecase returns Forbidden)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newDeleteRouterWithUserID(uc, 3)

		in := user_usecase.DeleteUserInput{EditorID: 3, TargetID: 4}
		uc.On("Delete", mock.Anything, in).Return(apperror.Forbiddenf("forbidden", nil))

		w := performDELETE(r, "/users/4")

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.JSONEq(t, `{"error":"forbidden"}`, w.Body.String())
		uc.AssertExpectations(t)
	})

	t.Run("500 - internal (usecase returns Internal)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		r := newDeleteRouterWithUserID(uc, 10)

		in := user_usecase.DeleteUserInput{EditorID: 10, TargetID: 11}
		uc.On("Delete", mock.Anything, in).Return(apperror.Internalf("internal error", nil))

		w := performDELETE(r, "/users/11")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"internal error"}`, w.Body.String())
		uc.AssertExpectations(t)
	})
}
