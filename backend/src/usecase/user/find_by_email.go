// app/usecase/user/find_by_email.go
package user

import (
	"context"
)

type FindByEmailInput struct {
	Email string
}

type FindByEmailOutput struct {
	UserID         int
	HashedPassword string
	IsAdmin        bool
	IsRoot         bool
	IsTest         bool
	// 必要ならEmailやNameも返せる
}

func (uc *UserUsecase) FindByEmail(ctx context.Context, in FindByEmailInput) (*FindByEmailOutput, error) {
	u, err := uc.userRepo.FindActiveByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	return &FindByEmailOutput{
		UserID:         u.ID,
		HashedPassword: u.Password,
		IsAdmin:        u.IsAdmin,
		IsRoot:         u.IsRoot,
		IsTest:         u.IsTest,
	}, nil
}
