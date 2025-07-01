package settingUc

import (
	"context"

	"word_app/backend/src/domain"
	settingport "word_app/backend/src/interfaces/repository/setting"
	"word_app/backend/src/interfaces/repository/tx"
)

type UpdateUserConfigInput struct {
	UserID     int  // 取得済み (Auth MW 等で)
	IsDarkMode bool // true=ダーク・false=ライト
}

type UpdateUserConfig struct {
	Repo settingport.UserConfigRepository
	Tx   tx.TxManager // 既存の Tx ラッパーを流用
}

func (uc *UpdateUserConfig) UpdateUserConfigExecute(ctx context.Context, in UpdateUserConfigInput) (*domain.UserConfig, error) {
	var out *domain.UserConfig
	err := uc.Tx.WithTx(ctx, func(txCtx context.Context) error {
		cfg := &domain.UserConfig{UserID: in.UserID, IsDarkMode: in.IsDarkMode}
		var err error
		out, err = uc.Repo.Upsert(txCtx, cfg)
		return err
	})
	return out, err
}
