package setting_service

import (
	"context"
	"word_app/backend/ent"
	"word_app/backend/ent/rootconfig"
	"word_app/backend/ent/user"

	"github.com/sirupsen/logrus"
)

func (e *EntSettingClient) GetRootConfig(ctx context.Context, userID int) (*ent.RootConfig, error) {
	// 管理者チェック
	userEntity, err := e.client.User().
		Query().
		Where(user.ID(userID)).
		Only(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, ErrDatabaseFailure
	}
	if !userEntity.IsRoot {
		return nil, ErrUnauthorized
	}

	rootConfig, err := e.client.RootConfig().
		Query().
		Where(rootconfig.ID(1)).
		First(ctx)

	if err != nil {
		return nil, ErrRootConfigNotFound
	}
	return rootConfig, nil
}
