package settingUc

import (
	"context"
	"errors"

	"word_app/backend/src/domain"
	settingRepo "word_app/backend/src/interfaces/repository/setting"
	userRepo "word_app/backend/src/interfaces/repository/user"
)

var (
	ErrRootConfigMissing = errors.New("root-config not found")
)

type GetRootConfigInput struct {
	UserID int
}

type GetRootConfigOutput struct {
	Config *domain.RootConfig
}

type GetRootConfigUsecase struct {
	userRepo userRepo.UserRepository
	cfgRepo  settingRepo.RootConfigRepository
}

func NewGetRootConfigUsecase(u userRepo.UserRepository, c settingRepo.RootConfigRepository) *GetRootConfigUsecase {
	return &GetRootConfigUsecase{userRepo: u, cfgRepo: c}
}

func (uc *GetRootConfigUsecase) GetRootConfigExecute(ctx context.Context, in GetRootConfigInput) (*GetRootConfigOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, in.UserID)
	if err != nil {
		return nil, err // ← DB エラーなどはそのまま
	}
	if !user.IsRoot {
		return nil, ErrUnauthorized
	}

	cfg, err := uc.cfgRepo.Get(ctx)
	if err != nil {
		return nil, ErrRootConfigMissing
	}
	return &GetRootConfigOutput{Config: cfg}, nil
}
