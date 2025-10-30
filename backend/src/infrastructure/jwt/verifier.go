// infrastructure/jwt/verifier.go
package jwt

import (
	"context"
	"errors"
	"fmt"

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
			// ベース非nil + wrap
			return nil, fmt.Errorf("%w", ErrUnexpectedAlg)
			// もし repoerr を使いたいなら:
			// return nil, repoerr.FromEnt(ErrUnexpectedAlg, ErrUnexpectedAlg.Error(), "")
		}
		return v.secret, nil
	})
	if err != nil || !tok.Valid {
		return "", fmt.Errorf("%w", ErrTokenInvalid)
		// or: return "", repoerr.FromEnt(ErrTokenInvalid, ErrTokenInvalid.Error(), "")
	}

	c, ok := tok.Claims.(*Claims)
	if !ok || c.UserID == "" {
		return "", fmt.Errorf("%w", ErrClaimsInvalid)
		// or: return "", repoerr.FromEnt(ErrClaimsInvalid, ErrClaimsInvalid.Error(), "")
	}
	return c.UserID, nil
}

var (
	ErrUnexpectedAlg = errors.New("unexpected signing method")
	ErrTokenInvalid  = errors.New("token not valid")
	ErrClaimsInvalid = errors.New("claims invalid")
)
