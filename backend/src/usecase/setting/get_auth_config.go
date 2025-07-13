package settingUc

import (
	"context"

	settingport "word_app/backend/src/infrastructure/repository/setting"
)

type AuthConfigDTO struct {
	IsLineAuth bool `json:"is_line_auth"`
}

type AuthConfigInteractor struct {
	repo settingport.RootConfigRepository
}

type GetAuthConfig interface {
	Execute(ctx context.Context) (*AuthConfigDTO, error)
}

func NewAuthConfig(r settingport.RootConfigRepository) *AuthConfigInteractor {
	return &AuthConfigInteractor{repo: r}
}

func (u *AuthConfigInteractor) Execute(ctx context.Context) (*AuthConfigDTO, error) {
	cfg, err := u.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &AuthConfigDTO{IsLineAuth: cfg.IsLineAuthentication}, nil
}
