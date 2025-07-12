package auth

import (
	"context"
	"time"
	auth_repo "word_app/backend/src/infrastructure/repository/auth"
	user_repo "word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/interfaces/http/auth"
	"word_app/backend/src/utils/tempjwt"

	"github.com/coreos/go-oidc/v3/oidc"
)

type AuthUsecase struct {
	provider     AuthProvider
	userRepo     user_repo.UserRepository
	extAuthRepo  auth_repo.ExternalAuthRepository
	jwtGenerator auth.JWTGenerator
	tempJwtGen   TempTokenGenerator
}

func NewAuthUsecase(
	provider AuthProvider,
	userRepo user_repo.UserRepository,
	extAuthRepo auth_repo.ExternalAuthRepository,
	jwtGen auth.JWTGenerator,
	tempJwtGen TempTokenGenerator,
) *AuthUsecase {
	return &AuthUsecase{
		provider:     provider,
		userRepo:     userRepo,
		extAuthRepo:  extAuthRepo,
		jwtGenerator: jwtGen,
		tempJwtGen:   tempJwtGen,
	}
}

type AuthProvider interface {
	AuthURL(state, nonce string) string
	Exchange(ctx context.Context, code string) (*tempjwt.Identity, error)
	ValidateNonce(idTok *oidc.IDToken, expected string) error
}

type TempTokenGenerator interface {
	GenerateTemp(id *tempjwt.Identity, ttl time.Duration) (string, error)
	ParseTemp(tok string) (*tempjwt.Identity, error)
}
