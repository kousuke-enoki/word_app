package user_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"word_app/backend/ent"
	"word_app/backend/src/handlers/user/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestSignUpHandler 正常系・異常系のテストケース
func TestSignUpHandler(t *testing.T) {
	// モックのクライアントとユーザーを生成
	mockClient := new(mocks.ClientInterface)
	mockUser := new(mocks.UserClientInterface) // Userメソッドが返すインターフェースのモック

	// User メソッドの返り値として mockUser を返すように設定
	mockClient.On("User").Return(mockUser)

	// mockUserに対して動作を定義
	mockUser.On("Create").Return(mockUser) // Createメソッドは mockUser 自身を返す
	mockUser.On("SetEmail", mock.Anything).Return(mockUser)
	mockUser.On("SetName", mock.Anything).Return(mockUser)
	mockUser.On("SetPassword", mock.Anything).Return(mockUser)
	mockUser.On("Save", mock.Anything).Return(&ent.User{ID: 1, Email: "test@example.com"}, nil)

	// テストハンドラーを呼び出し
	router := setupRouter(mockClient)
	t.Run("Valid SignUp Request", func(t *testing.T) {
		signUpReqBody := `{"email": "test@example.com", "name": "Test User", "password": "password123"}`
		w := performSignUpRequest(t, router, signUpReqBody)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	t.Run("Missing Fields", func(t *testing.T) {
		signUpReqBody := `{"email": "test@example.com", "name": ""}`
		w := performSignUpRequest(t, router, signUpReqBody)

		// ステータスコードが400であることを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request")
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		invalidReqBody := `{"email": "test@example.com", "name": "Test User", "password": 12345}` // passwordが数値
		w := performSignUpRequest(t, router, invalidReqBody)

		// ステータスコードが400であることを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request")
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		// サーバー内エラーを引き起こすために、バインディングで利用されるエンティティを変更する
		signUpReqBody := `{"email": "test@example.com", "name": "Test User", "password": "password123"}`
		w := performSignUpRequest(t, router, signUpReqBody)

		// ステータスコードが500であることを確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "sign up failed")
	})
}

func setupRouter(mockClient *mocks.ClientInterface) {
	panic("unimplemented")
}

func performSignUpRequest(_ *testing.T, router *gin.Engine, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/users/sign_up", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
