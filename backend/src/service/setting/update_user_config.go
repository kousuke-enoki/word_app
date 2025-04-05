package setting_service

import (
	"context"
	"word_app/backend/ent"
	"word_app/backend/ent/userconfig"
)

func (e *EntSettingClient) UpdateUserConfig(ctx context.Context, userID int, isLightMode bool) (*ent.UserConfig, error) {
	// 既存チェック
	exists, err := e.client.UserConfig().
		Query().
		Where(userconfig.UserIDEQ(userID)).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		_, err = e.client.UserConfig().
			Update().
			Where(userconfig.UserIDEQ(userID)).
			SetIsDarkMode(isLightMode).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		return e.client.UserConfig().
			Query().
			Where(userconfig.UserIDEQ(userID)).
			Only(ctx)
	}
	// 新規作成
	return e.client.UserConfig().
		Create().
		SetUserID(userID).
		SetIsDarkMode(isLightMode).
		Save(ctx)
}
