// interfaces/client.go
package interfaces

import (
	"word_app/backend/ent"
	"word_app/backend/src/interfaces/service_interfaces"
)

// ClientInterface 全体を統一するクライアントインターフェース
type ClientInterface interface {
	UserClient
	WordService
	EntClient() *ent.Client // 必要なら直接Entクライアントも提供
	service_interfaces.EntClientInterface
}
