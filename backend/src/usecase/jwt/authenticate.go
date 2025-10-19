// usecase/auth/authenticate.go
package jwt

import (
	"context"
	"strconv"

	"word_app/backend/src/infrastructure/jwt"
	"word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/shared/ucerr"
)

type JwtUsecase struct {
	verifier jwt.TokenVerifier
	users    user.Repository
}

func NewJwtUsecase(v jwt.TokenVerifier, u user.Repository) *JwtUsecase {
	return &JwtUsecase{verifier: v, users: u}
}

type Authenticator interface {
	Authenticate(ctx context.Context, raw string) (models.Principal, error)
}

func (s *JwtUsecase) Authenticate(ctx context.Context, raw string) (models.Principal, error) {
	sub, err := s.verifier.VerifyAndExtractSubject(ctx, raw)
	if err != nil {
		return models.Principal{}, err
	}
	id, err := strconv.Atoi(sub)
	if err != nil {
		return models.Principal{}, ucerr.Unauthorized("Unauthorized")
	}

	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		return models.Principal{}, err
	}
	// ここで“アプリ規則”を集中させる
	if u.DeletedAt != nil { // RequireActiveUser 相当
		return models.Principal{}, ucerr.Unauthorized("Unauthorized")
	}

	return models.Principal{
		UserID:  id,
		IsAdmin: u.IsAdmin,
		IsRoot:  u.IsRoot,
		IsTest:  u.IsTest,
	}, nil
}
