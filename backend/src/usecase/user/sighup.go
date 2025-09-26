// app/usecase/user/signup.go
package user

import (
	"context"

	"word_app/backend/src/domain"
)

type SignUpInput struct {
	Email    string
	Name     string
	Password string
}

type SignUpOutput struct {
	UserID int
}

func (uc *UserUsecase) SignUp(ctx context.Context, in SignUpInput) (*SignUpOutput, error) {
	u, err := domain.NewUser(in.Name, &in.Email, &in.Password)
	if err != nil {
		return nil, err
	}

	// Tx開始（既存Txがあれば join される実装が理想）
	txCtx, done, err := uc.txm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	commit := false
	defer func() { _ = done(commit) }()

	createdUser, err := uc.userRepo.Create(txCtx, u)
	if err != nil {
		return nil, err
	}
	if err := uc.settingRepo.CreateDefault(txCtx, createdUser.ID); err != nil {
		return nil, err
	}

	commit = true
	if err := done(commit); err != nil {
		return nil, err
	}

	return &SignUpOutput{UserID: createdUser.ID}, nil
}
