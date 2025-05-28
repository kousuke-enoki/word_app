package interfaces

import (
	"context"
	"word_app/backend/src/usecase/auth"
)

// JWTGenerator インターフェースを定義
type JWTGenerator interface {
	GenerateJWT(userID string) (string, error)
}
type AuthClient interface {
	GenerateJWT(userID string) (string, error)
	StartLogin(ctx context.Context, state, nonce string) string
	HandleCallback(ctx context.Context, code, state, nonce string) (*auth.CallbackResult, error)
	CompleteSignUp(ctx context.Context, tempToken, pass string) (string, error)
}

// type JwtMiddleware interface {
// 	JwtCheckMiddleware() gin.HandlerFunc
// }
