// src/infrastructure/auth/jwt_validator.go
package jwt

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"

	"word_app/backend/ent/user"
	"word_app/backend/src/infrastructure/auth"
	"word_app/backend/src/interfaces/service_interfaces"
	"word_app/backend/src/utils/contextutil"
)

type JWTValidator struct {
	secret []byte
	client service_interfaces.EntClientInterface
}

func NewJWTValidator(secret string, client service_interfaces.EntClientInterface) *JWTValidator {
	return &JWTValidator{
		secret: []byte(secret),
		client: client,
	}
}

// Validate は JWT を検証し、DB からロールを取得して返す。
// エラー種別を細かく返すことでハンドラ側でハンドリングしやすい。
func (v *JWTValidator) Validate(ctx context.Context, tokenStr string) (contextutil.UserRoles, error) {
	var zero contextutil.UserRoles

	tok, err := jwt.ParseWithClaims(tokenStr, auth.Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return v.secret, nil
	})
	if err != nil || !tok.Valid {
		return zero, errors.New("token_invalid")
	}

	claims, ok := tok.Claims.(*auth.Claims)
	if !ok || claims.UserID == "" {
		return zero, errors.New("claims_invalid")
	}

	id, err := strconv.Atoi(claims.UserID)
	if err != nil {
		return zero, errors.New("user_id_parse_error")
	}

	u, err := v.client.User().
		Query().
		Where(user.ID(id)).
		Select(user.FieldIsAdmin, user.FieldIsRoot).
		Only(ctx)
	if err != nil {
		return zero, errors.New("user_not_found")
	}

	return contextutil.UserRoles{
		UserID:  id,
		IsAdmin: u.IsAdmin,
		IsRoot:  u.IsRoot,
	}, nil
}
