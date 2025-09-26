package auth

import (
	"context"
	"time"

	auth_repo "word_app/backend/src/infrastructure/repository/auth"
	setting_repo "word_app/backend/src/infrastructure/repository/setting"
	tx_repo "word_app/backend/src/infrastructure/repository/tx"
	user_repo "word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/interfaces/http/auth"
	"word_app/backend/src/utils/tempjwt"

	"github.com/coreos/go-oidc/v3/oidc"
)

type Usecase struct {
	txm          tx_repo.Manager
	provider     Provider
	userRepo     user_repo.Repository
	settingRepo  setting_repo.UserConfigRepository
	extAuthRepo  auth_repo.ExternalAuthRepository
	jwtGenerator auth.JWTGenerator
	tempJwtGen   TempTokenGenerator
}

func NewUsecase(
	txm tx_repo.Manager,
	provider Provider,
	userRepo user_repo.Repository,
	settingRepo setting_repo.UserConfigRepository,
	extAuthRepo auth_repo.ExternalAuthRepository,
	jwtGen auth.JWTGenerator,
	tempJwtGen TempTokenGenerator,
) *Usecase {
	return &Usecase{
		txm:          txm,
		provider:     provider,
		userRepo:     userRepo,
		settingRepo:  settingRepo,
		extAuthRepo:  extAuthRepo,
		jwtGenerator: jwtGen,
		tempJwtGen:   tempJwtGen,
	}
}

type Provider interface {
	AuthURL(state, nonce string) string
	Exchange(ctx context.Context, code string) (*tempjwt.Identity, error)
	ValidateNonce(idTok *oidc.IDToken, expected string) error
}

type TempTokenGenerator interface {
	GenerateTemp(id *tempjwt.Identity, ttl time.Duration) (string, error)
	ParseTemp(tok string) (*tempjwt.Identity, error)
}
