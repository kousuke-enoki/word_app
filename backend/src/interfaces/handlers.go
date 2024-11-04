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

type WordHandler interface {
	AllWordListHandler() gin.HandlerFunc
	WordShowHandler() gin.HandlerFunc
}

type UserClient interface {
	CreateUser(ctx context.Context, email, name, password string) (*ent.User, error)
}
