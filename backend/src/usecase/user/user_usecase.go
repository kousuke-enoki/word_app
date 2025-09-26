// usecase/user_detail.go
package user

import (
	"word_app/backend/src/infrastructure/repository/auth"
	"word_app/backend/src/infrastructure/repository/setting"
	"word_app/backend/src/infrastructure/repository/tx"
	"word_app/backend/src/infrastructure/repository/user"
)

type UserUsecase struct {
	txm         tx.Manager // Begin(ctx) (txCtx, done, err) も提供
	userRepo    user.Repository
	settingRepo setting.UserConfigRepository
	authRepo    auth.ExternalAuthRepository // SoftDeleteByUserID （例: LINE/OIDCなど）
}

func NewUserUsecase(
	txm tx.Manager,
	userRepo user.Repository,
	settingRepo setting.UserConfigRepository,
	authRepo auth.ExternalAuthRepository,
) *UserUsecase {
	return &UserUsecase{
		txm:         txm,
		userRepo:    userRepo,
		settingRepo: settingRepo,
		authRepo:    authRepo,
	}
}
