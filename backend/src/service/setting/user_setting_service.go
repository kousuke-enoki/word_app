package setting_service

import (
	"errors"
	"word_app/backend/src/interfaces/service_interfaces"
)

// ent.Client をラップして settingClient インターフェースを実装
type EntSettingClient struct {
	client service_interfaces.EntClientInterface
}

func NewEntSettingClient(client service_interfaces.EntClientInterface) *EntSettingClient {
	return &EntSettingClient{client: client}
}

var (
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicateID        = errors.New("duplicate ID")
	ErrDatabaseFailure    = errors.New("database failure")
	ErrUnauthorized       = errors.New("unauthorized error")
	ErrUserConfigNotFound = errors.New("user setting not found")
	ErrRootConfigNotFound = errors.New("root setting not found")
)
