package settinguc

import (
	"context"

	"word_app/backend/src/domain"
	settingRepo "word_app/backend/src/infrastructure/repository/setting"

	"github.com/sirupsen/logrus"
)

type InputGetUserConfig struct {
	UserID int
}

type OutputGetUserConfig struct {
	Config *domain.UserConfig
}

type GetUserConfigInteractor struct {
	repo settingRepo.UserConfigRepository
}

type GetUserConfig interface {
	Execute(ctx context.Context, in InputGetUserConfig) (*OutputGetUserConfig, error)
}

func NewGetUserConfig(r settingRepo.UserConfigRepository) *GetUserConfigInteractor {
	return &GetUserConfigInteractor{repo: r}
}

func (uc *GetUserConfigInteractor) Execute(ctx context.Context, in InputGetUserConfig) (*OutputGetUserConfig, error) {
	cfg, err := uc.repo.GetByUserID(ctx, in.UserID)
	if err != nil {
		return nil, ErrUserConfigNotFound
	}
	logrus.Info(cfg, " retrieved successfully")
	return &OutputGetUserConfig{Config: cfg}, nil
}
