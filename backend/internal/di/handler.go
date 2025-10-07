// internal/di/handler.go
package di

import (
	"word_app/backend/config"
	AuthH "word_app/backend/src/handlers/auth"
	quizH "word_app/backend/src/handlers/quiz"
	resultH "word_app/backend/src/handlers/result"
	settingH "word_app/backend/src/handlers/setting"
	userH "word_app/backend/src/handlers/user"
	wordH "word_app/backend/src/handlers/word"
	"word_app/backend/src/infrastructure/jwt"
	"word_app/backend/src/interfaces"
	middleware_interface "word_app/backend/src/interfaces/http/middleware"
	"word_app/backend/src/interfaces/http/quiz"
	"word_app/backend/src/interfaces/http/result"
	"word_app/backend/src/interfaces/http/setting"
	"word_app/backend/src/interfaces/http/word"
	jwt_middleware "word_app/backend/src/middleware/jwt"

	quizSvc "word_app/backend/src/service/quiz"
	resultSvc "word_app/backend/src/service/result"
	wordSvc "word_app/backend/src/service/word"
)

type Handlers struct {
	JWTMiD  middleware_interface.Middleware // JWT ミドルウェアは Handler ではなく、インターフェースとして定義
	Auth    AuthH.Handler
	Setting setting.Handler
	User    userH.Handler
	Word    word.Handler
	Quiz    quiz.Handler
	Result  result.Handler
}

func NewHandlers(config *config.Config, uc *UseCases, client interfaces.ClientInterface) *Handlers {
	jwtGen := jwt.NewMyJWTGenerator(config.JWT.Secret)
	authClient := jwt.NewJWTValidator(config.JWT.Secret, client)
	// 既存のservice 層は “薄い Facade” として存続させる想定
	return &Handlers{
		JWTMiD:  jwt_middleware.NewMiddleware(authClient),
		Auth:    AuthH.NewHandler(uc.Auth, jwtGen),
		Setting: settingH.NewHandler(uc.Setting),
		User:    userH.NewHandler(uc.User, jwtGen),
		Word:    wordH.NewHandler(wordSvc.NewWordService(client)),
		Quiz:    quizH.NewHandler(quizSvc.NewService(client)),
		Result:  resultH.NewHandler(resultSvc.NewService(client)),
	}
}
