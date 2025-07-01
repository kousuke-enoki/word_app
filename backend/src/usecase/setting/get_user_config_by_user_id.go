package settingUc

import (
	"context"

	"word_app/backend/src/domain"
	settingRepo "word_app/backend/src/interfaces/repository/setting"
)

type GetUserConfigInput struct {
	UserID int
}

type GetUserConfigOutput struct {
	Config *domain.UserConfig
}

type GetUserConfigUsecase struct {
	repo settingRepo.UserConfigRepository
}

func NewGetUserConfigUsecase(r settingRepo.UserConfigRepository) *GetUserConfigUsecase {
	return &GetUserConfigUsecase{repo: r}
}

func (uc *GetUserConfigUsecase) GetUserConfigExecute(ctx context.Context, in GetUserConfigInput) (*GetUserConfigOutput, error) {
	cfg, err := uc.repo.GetByUserID(ctx, in.UserID)
	if err != nil {
		return nil, ErrUserConfigNotFound
	}
	return &GetUserConfigOutput{Config: cfg}, nil
}
