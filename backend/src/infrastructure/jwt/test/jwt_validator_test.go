package jwt_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"

	"word_app/backend/ent/enttest"
	jwti "word_app/backend/src/infrastructure/jwt"
	"word_app/backend/src/test"
	"word_app/backend/src/utils/contextutil"

	_ "github.com/mattn/go-sqlite3"
)

const secret = "test_secret"

// テスト用トークン生成
func makeToken(t *testing.T, uid string, exp time.Duration, key []byte) string {
	t.Helper()
	claims := &jwti.Claims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uid,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(key)
	require.NoError(t, err)
	return s
}

func TestTokenValidator_Validate(t *testing.T) {
	// ❶ in-memory SQLite
	ec := enttest.Open(t, "sqlite3", "file:memdb?mode=memory&cache=shared&_fk=1")
	defer ec.Close()

	// ❷ テーブルごとにセットアップ
	mustCreate := func(isAdmin, isRoot bool) {
		_, err := ec.User.Create().
			SetEmail("emaill@example.com").
			SetPassword("Password123$").
			SetIsAdmin(isAdmin).
			SetIsRoot(isRoot).
			Save(context.Background())
		require.NoError(t, err)
	}
	mustCreate(true, false) // 成功ケース用ユーザ

	validator := jwti.NewJWTValidator(secret, test.RealEntClient{Client: ec})
	ctx := context.Background()

	tests := []struct {
		name   string
		token  string
		expect contextutil.UserRoles
		errSub string
	}{
		{
			name:   "success",
			token:  makeToken(t, "1", time.Hour, []byte(secret)),
			expect: contextutil.UserRoles{UserID: 1, IsAdmin: true, IsRoot: false},
		},
		{
			name:   "signature invalid",
			token:  makeToken(t, "1", time.Hour, []byte("wrong")),
			errSub: "token_invalid",
		},
		{
			name:   "expired",
			token:  makeToken(t, "1", -time.Hour, []byte(secret)),
			errSub: "token_invalid parse_error",
		},
		{
			name:   "claims invalid",
			token:  "abc.def.ghi",
			errSub: "token_invalid",
		},
		{
			name:   "user_not_found", // User 99 を作らない
			token:  makeToken(t, "99", time.Hour, []byte(secret)),
			errSub: "user_not_found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := validator.Validate(ctx, tc.token)

			if tc.errSub != "" {
				require.ErrorContains(t, err, tc.errSub)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expect, got)
		})
	}
}
