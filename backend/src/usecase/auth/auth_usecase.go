package auth

import (
	"context"
	"fmt"
	"time"
	"word_app/backend/src/domain"
	"word_app/backend/src/interfaces"
	auth_port "word_app/backend/src/usecase/port/auth"
)

type AuthUsecase struct {
	provider     auth_port.AuthProvider
	userRepo     auth_port.UserRepository
	extAuthRepo  auth_port.ExternalAuthRepository
	jwtGenerator interfaces.JWTGenerator
	tempJwtGen   auth_port.TempTokenGenerator
}

func NewAuthUsecase(
	provider auth_port.AuthProvider,
	userRepo auth_port.UserRepository,
	extAuthRepo auth_port.ExternalAuthRepository,
	jwtGen interfaces.JWTGenerator,
	tempJwtGen auth_port.TempTokenGenerator,
) *AuthUsecase {
	return &AuthUsecase{
		provider:     provider,
		userRepo:     userRepo,
		extAuthRepo:  extAuthRepo,
		jwtGenerator: jwtGen,
		tempJwtGen:   tempJwtGen,
	}
}

func (u *AuthUsecase) StartLogin(ctx context.Context, state, nonce string) string {
	return u.provider.AuthURL(state, nonce)
}

func (u *AuthUsecase) HandleCallback(ctx context.Context, code, state, nonce string) (*interfaces.CallbackResult, error) {
	id, err := u.provider.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// ユーザ検索
	user, _ := u.userRepo.FindByProvider(ctx, id.Provider, id.Subject)
	if user != nil {
		token, _ := u.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", user.ID))
		return &interfaces.CallbackResult{Token: token}, nil
	}
	// 初回登録フロー
	temp, _ := u.tempJwtGen.GenerateTemp(id, 5*time.Minute)
	return &interfaces.CallbackResult{
		NeedPassword:  true,
		TempToken:     temp,
		SuggestedMail: id.Email,
	}, nil
}

func (u *AuthUsecase) CompleteSignUp(ctx context.Context, tempToken, pass string) (string, error) {
	id, err := u.tempJwtGen.ParseTemp(tempToken)
	if err != nil {
		return "", err
	}

	user, err := domain.NewUser(id.Email, id.Name, pass)
	if err != nil {
		return "", err
	}

	ext := domain.NewExternalAuth(0, id.Provider, id.Subject)

	if err := u.userRepo.Create(ctx, user, ext); err != nil {
		return "", err
	}
	return u.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", user.ID))
}
