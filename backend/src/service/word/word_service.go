package word

import (
	"errors"

	"word_app/backend/config"
	"word_app/backend/src/infrastructure/repository/tx"
	"word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/interfaces"
)

type ServiceImpl struct {
	client   interfaces.ClientInterface
	userRepo user.Repository
	limits   *config.LimitsCfg
	txm      tx.Manager
}

func NewWordService(
	client interfaces.ClientInterface,
	userRepo user.Repository,
	limits *config.LimitsCfg,
	txm tx.Manager,
) *ServiceImpl {
	return &ServiceImpl{
		client:   client,
		userRepo: userRepo,
		limits:   limits,
		txm:      txm,
	}
}

var (
	ErrWordNotFound       = errors.New("word not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrDeleteWord         = errors.New("failed to delete word")
	ErrDatabaseFailure    = errors.New("database failure")
	ErrWordExists         = errors.New("there is already a word with the same name")
	ErrCreateWord         = errors.New("failed to create word")
	ErrCreateWordInfo     = errors.New("failed to create word info")
	ErrCreateJapaneseMean = errors.New("failed to create japanese mean")
)
