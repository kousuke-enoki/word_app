// user/handler_test.go
package user_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/src/handlers/user"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/mocks"
	"word_app/backend/src/models"

	user_mocks "word_app/backend/src/mocks/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMyPageHandler(t *testing.T) {
	// テスト用のGinコンテキストとHTTPリクエストを準備
	gin.SetMode(gin.TestMode)
	mockClient := new(user_mocks.MockUsecase)
	mockJWTGen := &mocks.MockJwtGenerator{}

	userHandler := user.NewHandler(mockClient, mockJWTGen)

	// 正常時
	t.Run("Success", func(t *testing.T) {
		// モックの設定

		mockClient.
			On("GetMyDetail", mock.Anything, 1).
			Return(&models.UserDetail{
				ID:      1,
				Name:    "Test User",
				IsAdmin: true,
				IsRoot:  false,
				IsTest:  false,
			}, nil)

		// HTTPリクエストとレスポンスのセットアップ
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/mypage", nil)
		// コンテキストにユーザーIDを設定
		uid := 1
		p := models.Principal{
			UserID:  uid,
			IsAdmin: false,
			IsRoot:  false,
			IsTest:  false,
		}
		jwt.SetPrincipal(c, p)
		handler := userHandler.MyPageHandler()
		handler(c)

		// アサーション
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスボディの内容を確認（小文字の"admin"と"name"に変更）
		expected := `{"user":{"id":1,"name":"Test User","isAdmin":true,"isRoot":false,"isTest":false},"isLogin":true}`
		assert.JSONEq(t, expected, w.Body.String())

		mockClient.AssertExpectations(t)
	})

	// ユーザーIDが設定されていない場合
	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/mypage", nil)
		handler := userHandler.MyPageHandler()
		handler(c)

		// アサーション
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"unauthorized"}`, w.Body.String())
	})

	// データベースエラーが発生した場合
	t.Run("Database_error", func(t *testing.T) {
		mockClient.
			On("GetMyDetail", mock.Anything, 123).
			Return(nil, errors.New("some database error"))
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)

		c.Request = httptest.NewRequest(http.MethodGet, "/mypage", nil)
		uid := 123
		p := models.Principal{
			UserID:  uid,
			IsAdmin: false,
			IsRoot:  false,
			IsTest:  false,
		}
		jwt.SetPrincipal(c, p)
		handler := userHandler.MyPageHandler()
		handler(c)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.JSONEq(t, `{"error":"internal error"}`, recorder.Body.String())
		mockClient.AssertExpectations(t)
	})
}
