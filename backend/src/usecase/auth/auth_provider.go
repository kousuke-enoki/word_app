package auth

import (
	"context"
	"fmt"
	"time"

	"word_app/backend/src/domain"
)

func (u *AuthUsecase) StartLogin(_ context.Context, state, nonce string) string {
	return u.provider.AuthURL(state, nonce)
}

func (u *AuthUsecase) HandleCallback(ctx context.Context, code string) (*CallbackResult, error) {
	id, err := u.provider.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// ユーザ検索
	user, _ := u.userRepo.FindByProvider(ctx, id.Provider, id.Subject)
	if user != nil {
		token, _ := u.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", user.ID))
		return &CallbackResult{Token: token}, nil
	}
	// 初回登録フロー
	temp, _ := u.tempJwtGen.GenerateTemp(id, 5*time.Minute)
	return &CallbackResult{
		NeedPassword:  true,
		TempToken:     temp,
		SuggestedMail: id.Email,
	}, nil
}

func (u *AuthUsecase) CompleteSignUp(ctx context.Context, tempToken string, pass *string) (string, error) {
	id, err := u.tempJwtGen.ParseTemp(tempToken)
	if err != nil {
		return "", err
	}

	user, err := domain.NewUser(id.Name, id.Email, pass)
	if err != nil {
		return "", err
	}
	// Tx開始
	txCtx, done, err := u.txm.Begin(ctx)
	if err != nil {
		return "", err
	}
	commit := false
	defer func() { _ = done(commit) }()

	ext := domain.NewExternalAuth(0, id.Provider, id.Subject)

	createdUser, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return "", err
	}
	ext.UserID = createdUser.ID
	err = u.extAuthRepo.Create(ctx, ext)
	if err != nil {
		return "", err
	}

	if err := u.settingRepo.CreateDefault(txCtx, createdUser.ID); err != nil {
		return "", err
	}

	commit = true
	if err := done(commit); err != nil {
		return "", err
	}
	return u.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", user.ID))
}
