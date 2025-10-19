// infrastructure/jwt/verifier.go
package jwt

import (
	"context"
	"errors"

	"word_app/backend/src/infrastructure/repoerr"

	jwt "github.com/golang-jwt/jwt/v4"
)

type HS256Verifier struct{ secret []byte }

func NewHS256Verifier(secret string) *HS256Verifier {
	return &HS256Verifier{secret: []byte(secret)}
}

type TokenVerifier interface {
	// 署名OKか？期限切れか？だけ判断し、アプリにとって意味のある主語（sub/userID）を返す
	VerifyAndExtractSubject(ctx context.Context, raw string) (subject string, err error)
}

func (v *HS256Verifier) VerifyAndExtractSubject(ctx context.Context, raw string) (string, error) {
	tok, err := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, repoerr.FromEnt(errors.New("unexpected signing method"), "unexpected signing method", "")
		}
		return v.secret, nil
	})
	if err != nil || !tok.Valid {
		return "", repoerr.FromEnt(err, "token not valid", "")
	}

	c, ok := tok.Claims.(*Claims)
	if !ok || c.UserID == "" {
		return "", repoerr.FromEnt(err, "claims invalid", "")
	}
	return c.UserID, nil
}
