// infrastructure/repository/user_ent.go
package user

import (
	"context"
	"errors"
	"time"

	"word_app/backend/ent"
	"word_app/backend/src/domain"
	"word_app/backend/src/domain/repository"
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
	ListUsers(ctx context.Context, f repository.UserListFilter) (*repository.UserListResult, error)
	FindForUpdate(ctx context.Context, id int) (*domain.User, error)
	UpdatePartial(ctx context.Context, targetID int, f *repository.UserUpdateFields) (*domain.User, error)
	DeleteIfTest(ctx context.Context, id int) (deleted bool, err error)
	Exists(ctx context.Context, id int) (bool, error)
	IsTest(ctx context.Context, id int) (bool, error)
	LockByID(ctx context.Context, tx *ent.Tx, userID int) error
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
