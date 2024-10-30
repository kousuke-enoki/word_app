package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"word_app/backend/ent"
	"word_app/backend/ent/migrate"
	"word_app/backend/src/handlers/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSignUpHandler 正常系・異常系のテストケース
func TestSignUpHandler(t *testing.T) {
	// テスト用のデータベース接続
	client, err := ent.Open("postgres", "host=localhost port=5433 user=postgres dbname=db_test password=password sslmode=disable")
	require.NoError(t, err)
	defer client.Close() // テスト後にデータベース接続を閉じる

	// マイグレーションを実行
	err = client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true))
	require.NoError(t, err)

	// トランザクションを開始
	tx, err := client.Tx(context.Background())
	require.NoError(t, err)

	// defer でロールバックを確実に実行する
	defer func() {
		err = tx.Rollback()
		require.NoError(t, err)
	}()

	// Ginのテストモードを使用
	gin.SetMode(gin.TestMode)

	// Ginのルーターをセットアップ
	router := gin.Default()
	router.POST("/users/sign_up", user.SignUpHandler(client))

	t.Run("Valid SignUp Request", func(t *testing.T) {
		signUpReqBody := `{"email": "test@example.com", "name": "Test User", "password": "password123"}`
		w := performSignUpRequest(t, router, signUpReqBody)

		// ステータスコードが200であることを確認
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

func performSignUpRequest(_ *testing.T, router *gin.Engine, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/users/sign_up", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
