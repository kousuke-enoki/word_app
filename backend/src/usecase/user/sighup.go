// app/usecase/user/signup.go
package user

import (
	"context"

	"word_app/backend/src/domain"
	"word_app/backend/src/interfaces/http/user"
)

func (uc *UserUsecase) SignUp(ctx context.Context, in user.SignUpInput) (*user.SignUpOutput, error) {
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

	return &user.SignUpOutput{UserID: createdUser.ID}, nil
}
