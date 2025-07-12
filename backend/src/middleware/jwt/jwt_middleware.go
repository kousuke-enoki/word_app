package jwt

import (
	"word_app/backend/src/interfaces/http/middleware"
)

type JwtMiddleware struct {
	tokenValidator middleware.TokenValidator
}

func NewJwtMiddleware(tokenValidator middleware.TokenValidator) *JwtMiddleware {
	return &JwtMiddleware{
		tokenValidator: tokenValidator,
	}
}
