package auth

import (
	"context"
	"time"

	auth_repo "word_app/backend/src/infrastructure/repository/auth"
	user_repo "word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/interfaces/http/auth"
	"word_app/backend/src/utils/tempjwt"
)

type Usecase struct {
	provider     Provider
	userRepo     user_repo.Repository
	extAuthRepo  auth_repo.ExternalAuthRepository
	jwtGenerator auth.JWTGenerator
	tempJwtGen   TempTokenGenerator
}

func NewUsecase(
	provider Provider,
	userRepo user_repo.Repository,
	extAuthRepo auth_repo.ExternalAuthRepository,
	jwtGen auth.JWTGenerator,
	tempJwtGen TempTokenGenerator,
) *Usecase {
	return &Usecase{
		provider:     provider,
		userRepo:     userRepo,
		extAuthRepo:  extAuthRepo,
		jwtGenerator: jwtGen,
		tempJwtGen:   tempJwtGen,
	}
}

type Provider interface {
	AuthURL(state, nonce string) string
	Exchange(ctx context.Context, code string) (*tempjwt.Identity, error)
	ValidateNonce(idTok, expected string) error
}

type TempTokenGenerator interface {
	GenerateTemp(id *tempjwt.Identity, ttl time.Duration) (string, error)
	ParseTemp(tok string) (*tempjwt.Identity, error)
}
