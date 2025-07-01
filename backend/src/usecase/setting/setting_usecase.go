package settingUc

import (
	"context"
	"errors"

	"word_app/backend/src/domain"
	settingRepo "word_app/backend/src/interfaces/repository/setting"
	settingport "word_app/backend/src/interfaces/repository/setting"
	"word_app/backend/src/interfaces/repository/tx"
	userRepo "word_app/backend/src/interfaces/repository/user"
)

type SettingService struct {
	rootRepo    settingRepo.RootConfigRepository
	userRepo    userRepo.UserRepository
	userCfgRepo settingRepo.UserConfigRepository
	tx          tx.TxManager
}

func NewSettingService(r settingRepo.RootConfigRepository, u userRepo.UserRepository,
	uc settingRepo.UserConfigRepository, tx tx.TxManager) *SettingService {
	return &SettingService{rootRepo: r, userRepo: u, userCfgRepo: uc, tx: tx}
}

type AuthConfigDTO struct {
	IsLineAuth bool `json:"is_line_auth"`
}

type ConfigInteractor struct {
	repo settingport.RootConfigRepository
}

func NewConfigUsecase(r settingport.RootConfigRepository) *ConfigInteractor {
	return &ConfigInteractor{repo: r}
}

var (
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicateID        = errors.New("duplicate ID")
	ErrDatabaseFailure    = errors.New("database failure")
	ErrUnauthorized       = errors.New("unauthorized error")
	ErrUserConfigNotFound = errors.New("user setting not found")
	ErrRootConfigNotFound = errors.New("root setting not found")
)

type ConfigUsecase interface {
	GetAuthConfig(ctx context.Context) (*AuthConfigDTO, error)
	GetRootConfigExecute(ctx context.Context, in GetRootConfigInput) (*GetRootConfigOutput, error)
	GetUserConfigExecute(ctx context.Context, in GetUserConfigInput) (*GetUserConfigOutput, error)
	UpdateRootConfigExecute(ctx context.Context, in UpdateRootConfigInput) (*domain.RootConfig, error)
	UpdateUserConfigExecute(ctx context.Context, in UpdateUserConfigInput) (*domain.UserConfig, error)
}
