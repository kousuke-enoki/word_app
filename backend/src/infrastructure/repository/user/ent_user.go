// infrastructure/repository/user_ent.go
package user

import (
	"context"
	"errors"
	"time"

	"word_app/backend/src/domain"
	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

type EntUserRepo struct {
	client serviceinterfaces.EntClientInterface
}

type Repository interface {
	FindByProvider(ctx context.Context, provider, sub string) (*domain.User, error)
	Create(ctx context.Context, u *domain.User) (user *domain.User, err error)
	FindByID(ctx context.Context, id int) (*domain.User, error)
	IsRoot(ctx context.Context, userID int) (bool, error)
	FindDetailByID(ctx context.Context, id int) (*domain.User, error) // 詳細用：preload 付き
	SoftDeleteByID(ctx context.Context, id int, t time.Time) error
	FindActiveByEmail(ctx context.Context, email string) (*domain.User, error)
	ListUsers(ctx context.Context, f domain.UserListFilter) (*domain.UserListResult, error)
}

func NewEntUserRepo(c serviceinterfaces.EntClientInterface) *EntUserRepo {
	return &EntUserRepo{client: c}
}

var (
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrDuplicateID     = errors.New("duplicate ID")
	ErrDatabaseFailure = errors.New("database failure")
	ErrUserNotFound    = errors.New("user not found")
	ErrCreateConfig    = errors.New("create user config failure")
	ErrUnauthorized    = errors.New("unauthorized")
)
