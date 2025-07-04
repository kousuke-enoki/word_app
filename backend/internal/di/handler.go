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
	"word_app/backend/src/interfaces/http/auth"
	middleware_interface "word_app/backend/src/interfaces/http/middleware"
	"word_app/backend/src/interfaces/http/setting"
	jwt_middleware "word_app/backend/src/middleware/jwt"

	quizSvc "word_app/backend/src/service/quiz"
	resultSvc "word_app/backend/src/service/result"
	userSvc "word_app/backend/src/service/user"
	wordSvc "word_app/backend/src/service/word"
)

type Handlers struct {
	JWTMiD  middleware_interface.JwtMiddleware // JWT ミドルウェアは Handler ではなく、インターフェースとして定義
	Auth    auth.AuthHandler
	Setting setting.SettingHandler
	User    interfaces.UserHandler
	Word    interfaces.WordHandler
	Quiz    interfaces.QuizHandler
	Result  interfaces.ResultHandler
}

func NewHandlers(config *config.Config, uc *UseCases, client interfaces.ClientInterface) *Handlers {
	jwtGen := jwt.NewMyJWTGenerator(config.JWT.Secret)
	authClient := jwt.NewJWTValidator(config.JWT.Secret, client)
	// 既存のservice 層は “薄い Facade” として存続させる想定
	return &Handlers{
		JWTMiD:  jwt_middleware.NewJwtMiddleware(authClient),
		Auth:    AuthH.NewAuthHandler(uc.Auth, jwtGen),
		Setting: settingH.NewAuthSettingHandler(uc.Setting),
		User:    userH.NewUserHandler(userSvc.NewEntUserClient(client), jwtGen),
		Word:    wordH.NewWordHandler(wordSvc.NewWordService(client)),
		Quiz:    quizH.NewQuizHandler(quizSvc.NewQuizService(client)),
		Result:  resultH.NewResultHandler(resultSvc.NewResultService(client)),
	}
}
