package user

import (
	"context"
	"errors"

	"word_app/backend/src/domain"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"
)

// GetMyDetail: /users/me
func (uc *UserUsecase) GetMyDetail(ctx context.Context, viewerID int) (*models.UserDetail, error) {
	me, err := uc.userRepo.FindDetailByID(ctx, viewerID)
	if err != nil {
		return nil, apperror.New(apperror.NotFound, "notFound", nil)
	}
	return toDTO(me), nil
}

func (uc *UserUsecase) GetDetailByID(ctx context.Context, viewerID, targetID int) (*models.UserDetail, error) {
	viewer, err := uc.userRepo.FindByID(ctx, viewerID)
	if err != nil {
		return nil, apperror.New(apperror.Unauthorized, "unauthorized", err)
	}
	if !viewer.IsAdmin {
		return nil, apperror.New(apperror.Forbidden, "forbidden", nil)
	}
	target, err := uc.userRepo.FindDetailByID(ctx, targetID)
	if err != nil {
		return nil, apperror.New(apperror.NotFound, "user not found", err)
	}
	return toDTO(target), nil
}

// Domain → DTO（表現層向け整形はここ or Presenter）
func toDTO(u *domain.User) *models.UserDetail {
	return &models.UserDetail{
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

var ErrForbidden = errors.New("forbidden")
