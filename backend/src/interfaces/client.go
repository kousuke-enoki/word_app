package interfaces

import (
	"word_app/backend/ent"
	"word_app/backend/src/interfaces/http/middleware"
	"word_app/backend/src/interfaces/http/quiz"
	"word_app/backend/src/interfaces/http/user"
	"word_app/backend/src/interfaces/http/word"
	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
	settingUc "word_app/backend/src/usecase/setting"
)

// ClientInterface 全体を統一するクライアントインターフェース
type ClientInterface interface {
	user.Client
	word.Service
	quiz.Service
	middleware.TokenValidator
	EntClient() *ent.Client
	serviceinterfaces.EntClientInterface
	settingUc.SettingFacade
}
