package user

import (
	"context"

	"word_app/backend/src/domain"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/shared/ucerr"
)

// GetMyDetail: /users/me
func (uc *UserUsecase) GetMyDetail(ctx context.Context, viewerID int) (*models.UserDetail, error) {
	me, err := uc.userRepo.FindDetailByID(ctx, viewerID)
	if err != nil {
		return nil, err
	}
	return toDTO(me), nil
}

func (uc *UserUsecase) GetDetailByID(ctx context.Context, viewerID, targetID int) (*models.UserDetail, error) {
	viewer, err := uc.userRepo.FindByID(ctx, viewerID)
	if err != nil {
		return nil, err
	}
	if !viewer.IsAdmin {
		return nil, ucerr.Forbidden("forbidden")
	}
	target, err := uc.userRepo.FindDetailByID(ctx, targetID)
	if err != nil {
		return nil, err
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
