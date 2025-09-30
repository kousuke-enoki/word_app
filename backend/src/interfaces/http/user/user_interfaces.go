package user

import (
	"context"

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
	MeHandler() gin.HandlerFunc
	ShowHandler() gin.HandlerFunc
}

type Usecase interface {
	Delete(ctx context.Context, in DeleteUserInput) error
	FindByEmail(ctx context.Context, email string) (*FindByEmailOutput, error)
	UpdateUser(ctx context.Context, in UpdateUserInput) (*models.UserDetail, error)
	SignUp(ctx context.Context, in SignUpInput) (*SignUpOutput, error)
	ListUsers(ctx context.Context, in ListUsersInput) (*UserListResponse, error)
	GetMyDetail(ctx context.Context, viewerID int) (*models.UserDetail, error)
	GetDetailByID(ctx context.Context, viewerID, targetID int) (*models.UserDetail, error)
}

// type Validator interface {
// 	SignUp(SignUpRequest *models.SignUpRequest) []*models.FieldError
// 	SignIn(SignInRequest *models.SignInRequest) []*models.FieldError
// }
