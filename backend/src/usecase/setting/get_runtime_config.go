package settinguc

import (
	"context"
	"time"

	"word_app/backend/ent"
	settingRepo "word_app/backend/src/infrastructure/repository/setting"
	clock "word_app/backend/src/usecase/clock"
)

type RuntimeConfigDTO struct {
	IsTestUserMode       bool   `json:"is_test_user_mode"`
	IsLineAuthentication bool   `json:"is_line_authentication"`
	Version              string `json:"version"`
}

type RuntimeConfigInteractor struct {
	repo  settingRepo.RootConfigRepository
	clock clock.Clock
}

type GetRuntimeConfig interface {
	Execute(ctx context.Context) (*RuntimeConfigDTO, error)
}

func NewRuntimeConfig(
	r settingRepo.RootConfigRepository,
	clock clock.Clock,
) *RuntimeConfigInteractor {
	return &RuntimeConfigInteractor{
		repo:  r,
		clock: clock,
	}
}

func (u *RuntimeConfigInteractor) Execute(ctx context.Context) (*RuntimeConfigDTO, error) {
	cfg, err := u.repo.Get(ctx)
	if err != nil {
		// RootConfigが見つからない場合はデフォルト値を返す
		if ent.IsNotFound(err) {
			return &RuntimeConfigDTO{
				IsTestUserMode:       false,
				IsLineAuthentication: false,
				Version:              u.clock.Now().Format(time.DateTime),
			}, nil
		}
		return nil, err
	}
	// version は UpdatedAt を time.DateTime 形式で返す
	return &RuntimeConfigDTO{
		IsTestUserMode:       cfg.IsTestUserMode,
		IsLineAuthentication: cfg.IsLineAuthentication,
		Version:              cfg.UpdatedAt.Format(time.DateTime),
	}, nil
}
