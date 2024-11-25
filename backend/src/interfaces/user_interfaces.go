// interfaces/handlers.go
package interfaces

import (
	"context"

	"word_app/backend/ent"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	SignUpHandler() gin.HandlerFunc
	SignInHandler() gin.HandlerFunc
	MyPageHandler() gin.HandlerFunc
}

type UserClient interface {
	CreateUser(ctx context.Context, email, name, password string) (*ent.User, error)
	FindUserByEmail(ctx context.Context, email string) (*ent.User, error)
	FindUserByID(ctx context.Context, id int) (*ent.User, error)
}
