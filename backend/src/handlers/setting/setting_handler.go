// handlers/setting_handler.go
package setting

import (
	"word_app/backend/src/interfaces"
)

type SettingHandler struct {
	settingService interfaces.SettingClient
}

func NewSettingHandler(client interfaces.SettingClient) *SettingHandler {
	return &SettingHandler{
		settingService: client,
	}
}
