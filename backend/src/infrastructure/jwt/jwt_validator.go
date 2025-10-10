package jwt

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"

	"word_app/backend/ent/user"
	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
	"word_app/backend/src/utils/contextutil"
)

type JwtTokenValidator struct {
	secret []byte
	client serviceinterfaces.EntClientInterface
}

func NewJWTValidator(secret string, client serviceinterfaces.EntClientInterface) *JwtTokenValidator {
	return &JwtTokenValidator{
		secret: []byte(secret),
		client: client,
	}
}

type TokenValidator interface {
	// raw JWT を検証してユーザ権限を返す
	Validate(ctx context.Context, token string) (contextutil.UserRoles, error)
}

// Validate は JWT を検証し、DB からロールを取得して返す。
func (v *JwtTokenValidator) Validate(ctx context.Context, tokenStr string) (contextutil.UserRoles, error) {
	var zero contextutil.UserRoles
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return v.secret, nil
	})
	if err != nil {
		logrus.Warnf("parse error: %v", err)
		return zero, errors.New("token_invalid parse_error")
	}
	if !tok.Valid {
		logrus.Warn("token not valid")
		return zero, errors.New("token_invalid token_not_valid")
	}

	claims, ok := tok.Claims.(*Claims)
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
		Select(user.FieldIsAdmin, user.FieldIsRoot, user.FieldIsTest).
		Only(ctx)
	if err != nil {
		return zero, errors.New("user_not_found")
	}
	return contextutil.UserRoles{
		UserID:  id,
		IsAdmin: u.IsAdmin,
		IsRoot:  u.IsRoot,
		IsTest:  u.IsTest,
	}, nil
}
