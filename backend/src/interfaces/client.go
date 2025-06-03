package interfaces

import (
	"word_app/backend/ent"
	"word_app/backend/src/interfaces/http/middleware"
	"word_app/backend/src/interfaces/service_interfaces"
)

// ClientInterface 全体を統一するクライアントインターフェース
type ClientInterface interface {
	UserClient
	SettingClient
	WordService
	QuizService
	middleware.TokenValidator
	EntClient() *ent.Client
	service_interfaces.EntClientInterface
}
