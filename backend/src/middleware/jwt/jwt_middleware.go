package jwt

import (
	"word_app/backend/src/interfaces/http/middleware"
)

type Middleware struct {
	tokenValidator middleware.TokenValidator
}

func NewMiddleware(tokenValidator middleware.TokenValidator) *Middleware {
	return &Middleware{
		tokenValidator: tokenValidator,
	}
}
