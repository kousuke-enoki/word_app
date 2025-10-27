// infrastructure/jwt/verifier_test.go
package jwt_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	stdjwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"

	jwti "word_app/backend/src/infrastructure/jwt"
)

const secret = "test_secret"

func makeHS256Token(t *testing.T, uid string, expFromNow time.Duration, key []byte) string {
	t.Helper()
	claims := &jwti.Claims{
		UserID: uid,
		RegisteredClaims: stdjwt.RegisteredClaims{
			Subject:   uid,
			ExpiresAt: stdjwt.NewNumericDate(time.Now().Add(expFromNow)),
		},
	}
	token := stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, claims)
	s, err := token.SignedString(key)
	require.NoError(t, err)
	return s
}

func makeRS256Token(t *testing.T, uid string, expFromNow time.Duration, priv *rsa.PrivateKey) string {
	t.Helper()
	claims := &jwti.Claims{
		UserID: uid,
		RegisteredClaims: stdjwt.RegisteredClaims{
			Subject:   uid,
			ExpiresAt: stdjwt.NewNumericDate(time.Now().Add(expFromNow)),
		},
	}
	token := stdjwt.NewWithClaims(stdjwt.SigningMethodRS256, claims)
	s, err := token.SignedString(priv)
	require.NoError(t, err)
	return s
}

func TestHS256Verifier_VerifyAndExtractSubject(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	verifier := jwti.NewHS256Verifier(secret)
	rsaPriv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	tests := []struct {
		name     string
		rawToken string
		wantSub  string
		wantErr  error // ← 文字列ではなくエラー型で
	}{
		{
			name:     "success_valid_HS256",
			rawToken: makeHS256Token(t, "42", time.Hour, []byte(secret)),
			wantSub:  "42",
		},
		{
			name:     "invalid_signature_wrong_secret",
			rawToken: makeHS256Token(t, "42", time.Hour, []byte("wrong")),
			wantErr:  jwti.ErrTokenInvalid,
		},
		{
			name:     "expired_token",
			rawToken: makeHS256Token(t, "42", -1*time.Hour, []byte(secret)),
			wantErr:  jwti.ErrTokenInvalid,
		},
		{
			name:     "malformed_token_string",
			rawToken: "abc.def.ghi",
			wantErr:  jwti.ErrTokenInvalid,
		},
		{
			name:     "unexpected_signing_method_RS256",
			rawToken: makeRS256Token(t, "42", time.Hour, rsaPriv),
			wantErr:  jwti.ErrTokenInvalid,
		},
		{
			name:     "claims_invalid_empty_user_id",
			rawToken: makeHS256Token(t, "", time.Hour, []byte(secret)),
			wantErr:  jwti.ErrClaimsInvalid,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sub, err := verifier.VerifyAndExtractSubject(ctx, tc.rawToken)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantSub, sub)
		})
	}
}
