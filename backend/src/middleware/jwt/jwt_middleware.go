package jwt

import (
	"word_app/backend/src/infrastructure/jwt"

	"github.com/gin-gonic/gin"
)

type JwtMiddleware struct {
	tokenValidator jwt.TokenValidator
}

func NewMiddleware(tokenValidator jwt.TokenValidator) *JwtMiddleware {
	return &JwtMiddleware{
		tokenValidator: tokenValidator,
	}
}

type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
	JwtCheckMiddleware() gin.HandlerFunc
}
