package middleware

import (
	"context"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
)

type JwtMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
}
type TokenValidator interface {
	// raw JWT を検証してユーザ権限を返す
	Validate(ctx context.Context, token string) (contextutil.UserRoles, error)
}

type UserRoles struct {
	UserID  int
	IsAdmin bool
	IsRoot  bool
}
