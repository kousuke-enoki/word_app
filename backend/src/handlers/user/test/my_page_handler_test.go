// user/handler_test.go
package user_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/src/handlers/user"
	"word_app/backend/src/mocks"
	"word_app/backend/src/models"

	user_mocks "word_app/backend/src/mocks/http/user"

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
		c.Set("userID", 1)
		handler := userHandler.MyPageHandler()
		handler(c)

		// アサーション
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスボディの内容を確認（小文字の"admin"と"name"に変更）
		expected := `{"user":{"id":1,"name":"Test User","isAdmin":true,"isRoot":false,"isTest":false},"isLogin":true}`
		assert.JSONEq(t, expected, w.Body.String())

		mockClient.AssertExpectations(t)
	})

	// userIDの型がstringの時
	t.Run("type_error", func(t *testing.T) {
		// モックの設定
		userID := 1
		mockUser := &ent.User{
			ID:      1,
			Name:    "Test User",
			IsAdmin: true,
			IsRoot:  false,
			IsTest:  false,
		}
		mockClient.On("FindByID", mock.MatchedBy(func(c context.Context) bool {
			// コンテキストの確認
			return c != nil
		}), userID).Return(mockUser, nil)

		// HTTPリクエストとレスポンスのセットアップ
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/mypage", nil)
		// コンテキストにユーザーIDを設定
		c.Set("userID", "1")
		handler := userHandler.MyPageHandler()
		handler(c)

		// アサーション
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"unauthorized: userID not found in context"}`, w.Body.String())
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
		assert.JSONEq(t, `{"error":"unauthorized: userID not found in context"}`, w.Body.String())
	})

	// データベースエラーが発生した場合
	t.Run("Database_error", func(t *testing.T) {
		// Mock設定: FindByIDがエラーを返すようにする
		mockClient.
			On("GetMyDetail", mock.Anything, 123).
			Return(nil, errors.New("some database error"))
		// Ginのテストコンテキスト作成
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)

		c.Request = httptest.NewRequest(http.MethodGet, "/mypage", nil)
		// リクエスト設定
		req := httptest.NewRequest("GET", "/mypage", nil)
		c.Request = req

		// コンテキストに整数型のユーザーIDを設定
		c.Set("userID", 123)

		// ハンドラーを実行
		handler := userHandler.MyPageHandler()
		handler(c)

		// レスポンスを検証
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.JSONEq(t, `{"error":"internal error"}`, recorder.Body.String())
	})
}
