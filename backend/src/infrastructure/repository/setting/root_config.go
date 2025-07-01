package setting

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/rootconfig"
	"word_app/backend/src/domain"
	"word_app/backend/src/interfaces/service_interfaces"
)

type EntRootConfigRepo struct {
	client service_interfaces.EntClientInterface
}

func NewEntRootConfigRepo(c service_interfaces.EntClientInterface) *EntRootConfigRepo {
	return &EntRootConfigRepo{client: c}
}

func (r *EntRootConfigRepo) Get(ctx context.Context) (*domain.RootConfig, error) {
	rc, err := r.client.RootConfig().
		Query().
		Where(rootconfig.ID(1)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.RootConfig{
		ID:                   rc.ID,
		IsLineAuthentication: rc.IsLineAuthentication,
	}, nil
}

func (r *EntRootConfigRepo) Upsert(ctx context.Context, d *domain.RootConfig) (*domain.RootConfig, error) {
	// 存在確認
	rc, err := r.client.RootConfig().Query().Only(ctx)
	if ent.IsNotFound(err) {
		rc, err = r.client.RootConfig().
			Create().
			SetEditingPermission(d.EditingPermission).
			SetIsTestUserMode(d.IsTestUserMode).
			SetIsEmailAuthenticationCheck(d.IsEmailAuthenticationCheck).
			SetIsLineAuthentication(d.IsLineAuthentication).
			Save(ctx)
	} else if err == nil {
		rc, err = r.client.RootConfig().
			UpdateOne(rc).
			SetEditingPermission(d.EditingPermission).
			SetIsTestUserMode(d.IsTestUserMode).
			SetIsEmailAuthenticationCheck(d.IsEmailAuthenticationCheck).
			SetIsLineAuthentication(d.IsLineAuthentication).
			Save(ctx)
	}
	if err != nil {
		return nil, err
	}
	return entToDomain(rc), nil
}

func entToDomain(rc *ent.RootConfig) *domain.RootConfig {
	return &domain.RootConfig{
		ID:                         rc.ID,
		EditingPermission:          rc.EditingPermission,
		IsTestUserMode:             rc.IsTestUserMode,
		IsEmailAuthenticationCheck: rc.IsEmailAuthenticationCheck,
		IsLineAuthentication:       rc.IsLineAuthentication,
	}
}
