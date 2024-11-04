package user

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"word_app/backend/ent"
	"word_app/backend/ent/migrate"
	"word_app/backend/ent/user"
	"word_app/backend/src/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignInHandler(t *testing.T) {
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

	// トランザクション内でテスト用のユーザーを作成
	testUser, err := tx.User.Create().
		SetName("Test User").
		SetEmail("testuser@example.com").
		SetPassword("Password1234").
		Save(context.Background())
	require.NoError(t, err)

	// テスト用のJWTトークンを生成
	token, err := utils.GenerateJWT(strconv.Itoa(testUser.ID))
	require.NoError(t, err)

	// テストのHTTPリクエストを準備
	r := gin.Default()
	r.GET("/signin", SignInHandler(tx.Client()))

	req, _ := http.NewRequest(http.MethodGet, "/mypage", nil)
	// トークンをAuthorizationヘッダーに設定
	req.Header.Set("Authorization", "Bearer "+token)

	// レスポンスを記録
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	log.Println("test")
	log.Println(r)
	log.Println(w)
	log.Println(w.Body)
	log.Println(req)
	log.Println(t)
	// ステータスコードとレスポンスの検証
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test User")

	// トランザクション内でユーザーの確認
	fetchedUser, err := tx.User.Query().Where(user.ID(testUser.ID)).Only(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Test User", fetchedUser.Name)
}
