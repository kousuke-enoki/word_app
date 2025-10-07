package interfaces

import (
	"word_app/backend/ent"
	"word_app/backend/src/interfaces/http/quiz"
	"word_app/backend/src/interfaces/http/word"
	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

// ClientInterface 全体を統一するクライアントインターフェース
type ClientInterface interface {
	word.Service
	quiz.Service
	EntClient() *ent.Client
	serviceinterfaces.EntClientInterface
}
