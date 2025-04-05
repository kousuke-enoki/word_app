package setting_service

import (
	"context"
	"word_app/backend/ent"
	"word_app/backend/ent/userconfig"
)

func (e *EntSettingClient) GetUserConfig(ctx context.Context, userId int) (*ent.UserConfig, error) {
	userConfig, err := e.client.UserConfig().
		Query().
		Where(userconfig.UserID(userId)).
		First(ctx)

	if err != nil {
		return nil, ErrUserConfigNotFound
	}
	return userConfig, nil
}
