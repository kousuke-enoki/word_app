// internal/di/middleware.go
package di

import (
	jwtmw "word_app/backend/src/middleware/jwt"
)

type Middlewares struct {
	Auth jwtmw.Middleware // ← インターフェースで公開
}

func NewMiddlewares(uc *UseCases) *Middlewares {
	// uc.Jwt は usecase 側の Authenticator/JwtUsecase(インターフェース)想定
	return &Middlewares{
		Auth: jwtmw.NewMiddleware(uc.Jwt), // ← デリファレンスしない
	}
}
