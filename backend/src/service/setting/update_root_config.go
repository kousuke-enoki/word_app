package setting_service

import (
	"context"
	"word_app/backend/ent"
	"word_app/backend/ent/user"

	"github.com/sirupsen/logrus"
)

func (e *EntSettingClient) UpdateRootConfig(
	ctx context.Context,
	userID int,
	editingPermission string,
	isTestUserMode bool,
	isEmailAuth bool,
) (*ent.RootConfig, error) {
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

	// RootConfig の存在確認
	existing, err := e.client.RootConfig().
		Query().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			// 無い場合は新規作成
			return e.client.RootConfig().
				Create().
				SetEditingPermission(editingPermission).
				SetIsTestUserMode(isTestUserMode).
				SetIsEmailAuthentication(isEmailAuth).
				Save(ctx)
		}
		logrus.Error(err)
		return nil, ErrDatabaseFailure
	}

	// 存在する場合は更新
	updated, err := e.client.RootConfig().
		UpdateOne(existing).
		SetEditingPermission(editingPermission).
		SetIsTestUserMode(isTestUserMode).
		SetIsEmailAuthentication(isEmailAuth).
		Save(ctx)

	if err != nil {
		logrus.Error(err)
		return nil, ErrDatabaseFailure
	}
	return updated, nil
}
