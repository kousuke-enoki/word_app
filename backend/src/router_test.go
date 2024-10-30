package src

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"word_app/backend/ent"
	"word_app/backend/ent/migrate"
	"word_app/backend/src/utils"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupRouter(t *testing.T) {
	// テスト用のデータベース接続
	client, err := ent.Open("postgres", "host=localhost port=5433 user=postgres dbname=db_test password=password sslmode=disable")
	require.NoError(t, err)
	defer client.Close() // テスト後にデータベース接続を閉じる

	// マイグレーションを実行
	err = client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true))
	require.NoError(t, err)

	// Ginのテスト用の設定
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	SetupRouter(router, client)

	// ヘルスチェックエンドポイントのテスト
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"healthy"`)

	// サインアップエンドポイントのテスト
	signUpReq, _ := http.NewRequest(http.MethodPost, "/users/sign_up", nil)
	signUpW := httptest.NewRecorder()

	router.ServeHTTP(signUpW, signUpReq)
	assert.Equal(t, http.StatusBadRequest, signUpW.Code) // リクエストボディが空のため400エラーを期待

	// JWTトークン生成と保護されたエンドポイントのテスト
	// テスト用のユーザーを作成
	testUser, err := client.User.Create().
		SetName("Test User").
		SetEmail("testuser@example.com").
		SetPassword("Password1234").
		Save(context.Background())
	require.NoError(t, err)

	// テスト用のJWTトークンを生成
	token, err := utils.GenerateJWT(strconv.Itoa(testUser.ID))
	require.NoError(t, err)

	// 認証ミドルウェアを通過するエンドポイントのテスト
	protectedReq, _ := http.NewRequest(http.MethodGet, "/users/my_page", nil)
	protectedReq.Header.Set("Authorization", "Bearer "+token)
	protectedW := httptest.NewRecorder()
	router.ServeHTTP(protectedW, protectedReq)
	assert.Equal(t, http.StatusOK, protectedW.Code)
	assert.Contains(t, protectedW.Body.String(), "Test User")

	// CORSミドルウェアのテスト
	corsReq, _ := http.NewRequest(http.MethodOptions, "/", nil)
	corsReq.Header.Set("Origin", "http://localhost:3000")
	corsW := httptest.NewRecorder()
	router.ServeHTTP(corsW, corsReq)
	assert.Equal(t, http.StatusNoContent, corsW.Code)
	assert.Equal(t, "http://localhost:3000", corsW.Header().Get("Access-Control-Allow-Origin"))
}
