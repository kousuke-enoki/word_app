package jwt

import "word_app/backend/src/interfaces"

type JwtMiddleware struct {
	tokenValidator interfaces.TokenValidator
}

func NewJwtMiddleware(tokenValidator interfaces.TokenValidator) *JwtMiddleware {
	return &JwtMiddleware{tokenValidator: tokenValidator}
}
