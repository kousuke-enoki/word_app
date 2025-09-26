package settinguc

import (
	"context"

	"word_app/backend/src/domain"
	settingRepo "word_app/backend/src/infrastructure/repository/setting"
	"word_app/backend/src/infrastructure/repository/tx"
)

type InputUpdateUserConfig struct {
	UserID     int  `json:"user_id"`      // 取得済み (Auth MW 等で)
	IsDarkMode bool `json:"is_dark_mode"` // true=ダーク・false=ライト
}

type UpdateUserConfigInteractor struct {
	Tx       tx.Manager // 既存の Tx ラッパーを流用
	userRepo settingRepo.UserConfigRepository
}

type UpdateUserConfig interface {
	Execute(ctx context.Context, in InputUpdateUserConfig) (*domain.UserConfig, error)
}

func NewUpdateUserConfig(tx tx.Manager, u settingRepo.UserConfigRepository) *UpdateUserConfigInteractor {
	return &UpdateUserConfigInteractor{Tx: tx, userRepo: u}
}

func (uc *UpdateUserConfigInteractor) Execute(ctx context.Context, in InputUpdateUserConfig) (*domain.UserConfig, error) {
	var out *domain.UserConfig
	err := uc.Tx.WithTx(ctx, func(txCtx context.Context) error {
		cfg := &domain.UserConfig{UserID: in.UserID, IsDarkMode: in.IsDarkMode}
		var err error
		out, err = uc.userRepo.Upsert(txCtx, cfg) // ← txCtx を渡す
		return err
	})
	return out, err
}
