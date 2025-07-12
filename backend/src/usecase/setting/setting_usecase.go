package settingUc

import (
	"errors"
)

var (
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicateID        = errors.New("duplicate ID")
	ErrDatabaseFailure    = errors.New("database failure")
	ErrUnauthorized       = errors.New("unauthorized error")
	ErrUserConfigNotFound = errors.New("user setting not found")
	ErrRootConfigNotFound = errors.New("root setting not found")
)

// type ConfigUsecase interface {
// 	GetAuthConfig(ctx context.Context) (*AuthConfigDTO, error)
// 	GetRootConfigExecute(ctx context.Context, in GetRootConfigInput) (*GetRootConfigOutput, error)
// 	GetUserConfigExecute(ctx context.Context, in GetUserConfigInput) (*GetUserConfigOutput, error)
// 	UpdateRootConfigExecute(ctx context.Context, in UpdateRootConfigInput) (*domain.RootConfig, error)
// 	UpdateUserConfigExecute(ctx context.Context, in UpdateUserConfigInput) (*domain.UserConfig, error)
// }
