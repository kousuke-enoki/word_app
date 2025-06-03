package jwt

import (
	"word_app/backend/src/interfaces/http/middleware"
)

type JwtMiddleware struct {
	tokenValidator middleware.TokenValidator
}

// // JwtCheckMiddleware implements middleware.JwtMiddleware.
// func (j *JwtMiddleware) JwtCheckMiddleware() gin.HandlerFunc {
// 	panic("unimplemented")
// }

func NewJwtMiddleware(tokenValidator middleware.TokenValidator) *JwtMiddleware {
	return &JwtMiddleware{
		tokenValidator: tokenValidator,
	}
}
