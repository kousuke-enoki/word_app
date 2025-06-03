package jwt

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"

	"word_app/backend/ent/user"
	"word_app/backend/src/interfaces/service_interfaces"
	"word_app/backend/src/utils/contextutil"
)

type TokenValidator struct {
	secret []byte
	client service_interfaces.EntClientInterface
}

func NewJWTValidator(secret string, client service_interfaces.EntClientInterface) *TokenValidator {
	return &TokenValidator{
		secret: []byte(secret),
		client: client,
	}
}

// Validate は JWT を検証し、DB からロールを取得して返す。
// エラー種別を細かく返すことでハンドラ側でハンドリングしやすい。
func (v *TokenValidator) Validate(ctx context.Context, tokenStr string) (contextutil.UserRoles, error) {
	var zero contextutil.UserRoles
	logrus.Info("validate")
	logrus.Info(tokenStr)
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		logrus.Info("1")
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Info("2")
			return nil, fmt.Errorf("unexpected signing method")
		}
		logrus.Info("3")
		return v.secret, nil
	})
	if err != nil {
		logrus.Warnf("parse error: %v", err)
		logrus.Info("parse error")
		return zero, errors.New("token_invalid parse_error")
	}
	if !tok.Valid {
		logrus.Warn("token not valid")
		logrus.Info("token not valid")
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
		Select(user.FieldIsAdmin, user.FieldIsRoot).
		Only(ctx)
	if err != nil {
		return zero, errors.New("user_not_found")
	}
	logrus.Info(id)
	logrus.Info(u.IsAdmin)
	logrus.Info(u.IsRoot)
	return contextutil.UserRoles{
		UserID:  id,
		IsAdmin: u.IsAdmin,
		IsRoot:  u.IsRoot,
	}, nil
}
