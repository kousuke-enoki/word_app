// app/usecase/user/find_by_email.go
package user

import (
	"context"

	"word_app/backend/src/interfaces/http/user"
)

func (uc *UserUsecase) FindByEmail(ctx context.Context, email string) (*user.FindByEmailOutput, error) {
	u, err := uc.userRepo.FindActiveByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user.FindByEmailOutput{
		UserID:         u.ID,
		HashedPassword: u.Password,
		IsAdmin:        u.IsAdmin,
		IsRoot:         u.IsRoot,
		IsTest:         u.IsTest,
	}, nil
}
