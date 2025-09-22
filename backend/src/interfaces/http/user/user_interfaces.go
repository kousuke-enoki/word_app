package user

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	SignUpHandler() gin.HandlerFunc
	SignInHandler() gin.HandlerFunc
	MyPageHandler() gin.HandlerFunc
	ListHandler() gin.HandlerFunc
	EditHandler() gin.HandlerFunc
	DeleteHandler() gin.HandlerFunc
}

type Client interface {
	Create(ctx context.Context, email, name, password string) (*ent.User, error)
	FindByEmail(ctx context.Context, email string) (*ent.User, error)
	FindByID(ctx context.Context, id int) (*ent.User, error)
	GetUsers(ctx context.Context, UserListRequest *models.UserListRequest) (*models.UserListResponse, error)
	Update(ctx context.Context, UserEditRequest *models.UpdateUserInput) (*ent.User, error)
	Delete(ctx context.Context, editorID, targetID int) error
}

type Validator interface {
	SignUp(SignUpRequest *models.SignUpRequest) []*models.FieldError
	SignIn(SignInRequest *models.SignInRequest) []*models.FieldError
}
