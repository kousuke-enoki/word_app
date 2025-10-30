package jwt

import (
	"word_app/backend/src/usecase/jwt"

	"github.com/gin-gonic/gin"
)

type JwtMiddleware struct {
	JwtUsecase jwt.Authenticator
}

func NewMiddleware(JwtUsecase jwt.Authenticator) *JwtMiddleware {
	return &JwtMiddleware{
		JwtUsecase: JwtUsecase,
	}
}

type Middleware interface {
	AuthenticateMiddleware() gin.HandlerFunc
	// JwtCheckMiddleware() gin.HandlerFunc
}
