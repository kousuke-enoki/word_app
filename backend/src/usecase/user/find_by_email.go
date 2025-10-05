// app/usecase/user/find_by_email.go
package user

import (
	"context"
)

func (uc *UserUsecase) FindByEmail(ctx context.Context, email string) (*FindByEmailOutput, error) {
	u, err := uc.userRepo.FindActiveByEmail(ctx, email)
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
