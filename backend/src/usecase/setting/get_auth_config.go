package settingUc

import (
	"context"
)

func (u *ConfigInteractor) GetAuthConfig(ctx context.Context) (*AuthConfigDTO, error) {
	cfg, err := u.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &AuthConfigDTO{IsLineAuth: cfg.IsLineAuthentication}, nil
}
