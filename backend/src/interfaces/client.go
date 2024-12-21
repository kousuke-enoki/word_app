// interfaces/client.go
package interfaces

import (
	"word_app/backend/ent"
)

// ClientInterface 全体を統一するクライアントインターフェース
type ClientInterface interface {
	UserClient
	WordService
	EntClient() *ent.Client // 必要なら直接Entクライアントも提供
}
