// usecase/user_detail.go
package usecase

import (
	"context"
	"net/http"

	"word_app/backend/src/domain"
	user_repo "word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/models"
)

type UserDetailUsecase struct {
	Repo user_repo.EntUserRepo
}

func NewUserDetailUsecase(repo user_repo.EntUserRepo) *UserDetailUsecase {
	return &UserDetailUsecase{Repo: repo}
}

// GetMyDetail: /users/me
func (u *UserDetailUsecase) GetMyDetail(ctx context.Context, viewerID int) (*domain.User, error) {
	entUser, err := u.Repo.FindDetailByID(ctx, viewerID)
	if err != nil {
		return nil, err
	}
	return toUserDetail(entUser), nil
}

func (uc *UserDetailUsecase) GetMyDetail(ctx context.Context, viewerID int) (*models.User, int, error) {
	me, err := uc.Repo.FindDetailByID(ctx, viewerID)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	return toDTO(me), http.StatusOK, nil
}

func (uc *UserDetailUsecase) GetDetailByID(ctx context.Context, viewerID, targetID int) (*models.User, int, error) {
	viewer, err := uc.Repo.FindByID(ctx, viewerID)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	if !viewer.IsAdmin { // Admin 以上のみ
		return nil, http.StatusForbidden, ErrForbidden
	}
	target, err := uc.Repo.FindDetailByID(ctx, targetID)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return toDTO(target), http.StatusOK, nil
}

// Domain → DTO（表現層向け整形はここ or Presenter）
func toDTO(u *domain.User) *models.User {
	return &models.User{
		ID:               u.ID,
		Name:             u.Name,
		Email:            u.Email,
		IsAdmin:          u.IsAdmin,
		IsRoot:           u.IsRoot,
		IsTest:           u.IsTest,
		IsLine:           u.HasLine,
		IsSettedPassword: u.HasPassword,
		CreatedAt:        u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
