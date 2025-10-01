// user/handler_test.go
package user_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/src/handlers/user"
	"word_app/backend/src/mocks"

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
		userID := 1
		mockUser := &ent.User{
			Name:    "Test User",
			IsAdmin: true,
		}
		mockClient.On("FindByID", mock.Anything, userID).Return(mockUser, nil)

		// HTTPリクエストとレスポンスのセットアップ
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// コンテキストにユーザーIDを設定
		c.Set("userID", 1)
		handler := userHandler.MyPageHandler()
		handler(c)

		// アサーション
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスボディの内容を確認（小文字の"admin"と"name"に変更）
		expectedResponse := `{"user":{"id":0, "isAdmin":true, "isRoot":false,"isTest":false, "name":"Test User"}, "isLogin":true}`
		assert.JSONEq(t, expectedResponse, w.Body.String())

		mockClient.AssertExpectations(t)
	})

	// userIDの型がstringの時
	t.Run("type_error", func(t *testing.T) {
		// モックの設定
		userID := 1
		mockUser := &ent.User{
			Name:    "Test User",
			IsAdmin: true,
		}
		mockClient.On("FindByID", mock.MatchedBy(func(c context.Context) bool {
			// コンテキストの確認
			return c != nil
		}), userID).Return(mockUser, nil)

		// HTTPリクエストとレスポンスのセットアップ
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// コンテキストにユーザーIDを設定
		c.Set("userID", "1")
		handler := userHandler.MyPageHandler()
		handler(c)

		// アサーション
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// レスポンスボディの内容を確認（小文字の"admin"と"name"に変更）
		expectedResponse := `{"error":"Invalid userID type"}`
		assert.JSONEq(t, expectedResponse, w.Body.String())

		mockClient.AssertExpectations(t)
	})

	// ユーザーIDが設定されていない場合
	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler := userHandler.MyPageHandler()
		handler(c)

		// アサーション
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
		mockClient.AssertExpectations(t)
	})

	// データベースエラーが発生した場合
	t.Run("Database_error", func(t *testing.T) {
		// Mock設定: FindByIDがエラーを返すようにする
		mockClient.On("FindByID", mock.Anything, 123).
			Return(nil, errors.New("some database error"))
		// Ginのテストコンテキスト作成
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)

		// リクエスト設定
		req := httptest.NewRequest("GET", "/mypage", nil)
		c.Request = req

		// コンテキストに整数型のユーザーIDを設定
		c.Set("userID", 123)

		// ハンドラーを実行
		handler := userHandler.MyPageHandler()
		handler(c)

		// レスポンスを検証
		if status := recorder.Code; status != http.StatusInternalServerError {
			t.Errorf("expected status code 500, got %d", status)
		}
		if !strings.Contains(recorder.Body.String(), "Failed to retrieve user") {
			t.Errorf("expected error message 'Failed to retrieve user', got %s", recorder.Body.String())
		}

		// Mockが正しく呼び出されたかを確認
		mockClient.AssertExpectations(t)
	})
}
