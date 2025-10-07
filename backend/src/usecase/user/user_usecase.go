// usecase/user_usecase.go
package user

import (
	"context"

	"word_app/backend/src/infrastructure/repository/auth"
	"word_app/backend/src/infrastructure/repository/setting"
	"word_app/backend/src/infrastructure/repository/tx"
	"word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/models"
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

type Usecase interface {
	Delete(ctx context.Context, in DeleteUserInput) error
	FindByEmail(ctx context.Context, email string) (*FindByEmailOutput, error)
	UpdateUser(ctx context.Context, in UpdateUserInput) (*models.UserDetail, error)
	SignUp(ctx context.Context, in SignUpInput) (*SignUpOutput, error)
	ListUsers(ctx context.Context, in ListUsersInput) (*UserListResponse, error)
	GetMyDetail(ctx context.Context, viewerID int) (*models.UserDetail, error)
	GetDetailByID(ctx context.Context, viewerID, targetID int) (*models.UserDetail, error)
}
