// interfaces/client.go
package interfaces

import (
	"word_app/backend/ent"
	"word_app/backend/src/interfaces/service_interfaces"
)

// ClientInterface 全体を統一するクライアントインターフェース
type ClientInterface interface {
	UserClient
	SettingClient
	WordService
	EntClient() *ent.Client
	service_interfaces.EntClientInterface
}
