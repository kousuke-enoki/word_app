package setting_service

import (
	"context"
	"word_app/backend/ent/rootconfig"
	"word_app/backend/src/models"
)

func (e *EntSettingClient) GetAuthConfig(ctx context.Context) (*models.AuthSettingResponse, error) {
	rootConfig, err := e.client.RootConfig().
		Query().
		Where(rootconfig.ID(1)).
		First(ctx)

	if err != nil {
		return nil, ErrRootConfigNotFound
	}
	response := &models.AuthSettingResponse{
		IsLineAuth: rootConfig.IsLineAuthentication,
	}
	return response, nil
}
