package auth

import (
	"context"

	"word_app/backend/src/infrastructure/auth/line"
	"word_app/backend/src/infrastructure/jwt"
	auth_repo "word_app/backend/src/infrastructure/repository/auth"
	setting_repo "word_app/backend/src/infrastructure/repository/setting"
	tx_repo "word_app/backend/src/infrastructure/repository/tx"
	user_repo "word_app/backend/src/infrastructure/repository/user"
	userDailyUsageRepo "word_app/backend/src/infrastructure/repository/userdailyusage"
	clock "word_app/backend/src/usecase/clock"
)

type AuthUsecase struct {
	txm                tx_repo.Manager
	provider           line.Provider
	userRepo           user_repo.Repository
	settingRepo        setting_repo.UserConfigRepository
	extAuthRepo        auth_repo.ExternalAuthRepository
	jwtGenerator       jwt.JWTGenerator
	tempJwtGen         jwt.TempTokenGenerator
	rootSettingRepo    setting_repo.RootConfigRepository
	userDailyUsageRepo userDailyUsageRepo.Repository
	clock              clock.Clock
}

func NewUsecase(
	txm tx_repo.Manager,
	provider line.Provider,
	userRepo user_repo.Repository,
	settingRepo setting_repo.UserConfigRepository,
	extAuthRepo auth_repo.ExternalAuthRepository,
	jwtGen jwt.JWTGenerator,
	tempJwtGen jwt.TempTokenGenerator,
	rootSettingRepo setting_repo.RootConfigRepository,
	userDailyUsageRepo userDailyUsageRepo.Repository,
	clock clock.Clock,
) *AuthUsecase {
	return &AuthUsecase{
		txm:                txm,
		provider:           provider,
		userRepo:           userRepo,
		settingRepo:        settingRepo,
		extAuthRepo:        extAuthRepo,
		jwtGenerator:       jwtGen,
		tempJwtGen:         tempJwtGen,
		rootSettingRepo:    rootSettingRepo,
		userDailyUsageRepo: userDailyUsageRepo,
		clock:              clock,
	}
}

type Usecase interface {
	StartLogin(ctx context.Context, state, nonce string) string
	HandleCallback(ctx context.Context, code string) (*CallbackResult, error)
	CompleteSignUp(ctx context.Context, tempToken string, pass *string) (string, error)
	TestLogin(ctx context.Context) (*TestLoginOutput, error)
}

type CallbackResult struct {
	Token         string  `json:"token,omitempty"`
	NeedPassword  bool    `json:"need_password,omitempty"`
	TempToken     string  `json:"temp_token,omitempty"`
	SuggestedMail *string `json:"suggested_mail,omitempty"`
}
