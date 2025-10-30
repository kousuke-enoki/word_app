// internal/di/handler.go
package di

import (
	"word_app/backend/config"
	AuthH "word_app/backend/src/handlers/auth"
	BulkH "word_app/backend/src/handlers/bulk"
	quizH "word_app/backend/src/handlers/quiz"
	resultH "word_app/backend/src/handlers/result"
	settingH "word_app/backend/src/handlers/setting"
	userH "word_app/backend/src/handlers/user"
	wordH "word_app/backend/src/handlers/word"
	"word_app/backend/src/infrastructure/jwt"
	"word_app/backend/src/interfaces"
	"word_app/backend/src/interfaces/http/quiz"
	"word_app/backend/src/interfaces/http/result"
	"word_app/backend/src/interfaces/http/word"
)

type Handlers struct {
	Auth    AuthH.Handler
	Bulk    BulkH.Handler
	Setting settingH.Handler
	User    userH.Handler
	Word    word.Handler
	Quiz    quiz.Handler
	Result  result.Handler
}

func NewHandlers(config *config.Config, uc *UseCases, client interfaces.ClientInterface, s *Services) *Handlers {
	jwtGen := jwt.NewMyJWTGenerator(config.JWT.Secret)
	// 既存のservice 層は “薄い Facade” として存続させる想定
	return &Handlers{
		Auth:    AuthH.NewHandler(uc.Auth, jwtGen),
		Bulk:    BulkH.NewHandler(uc.BulkToken, uc.BulkRegister, &config.Limits),
		Setting: settingH.NewHandler(uc.Setting),
		User:    userH.NewHandler(uc.User, jwtGen),
		Word:    wordH.NewHandler(s.Word),
		Quiz:    quizH.NewHandler(s.Quiz),
		Result:  resultH.NewHandler(s.Result),
	}
}
