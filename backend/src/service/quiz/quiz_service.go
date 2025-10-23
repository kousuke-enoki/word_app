package quiz

import (
	"errors"

	"word_app/backend/config"
	"word_app/backend/src/infrastructure/repository/userdailyusage"
	"word_app/backend/src/interfaces"
	"word_app/backend/src/usecase/clock"
)

type ServiceImpl struct {
	client             interfaces.ClientInterface
	userDailyUsageRepo userdailyusage.Repository
	clock              clock.Clock
	limits             *config.LimitsCfg
}

func NewService(
	client interfaces.ClientInterface,
	userDailyUsageRepo userdailyusage.Repository,
	clock clock.Clock,
	limits *config.LimitsCfg,
) *ServiceImpl {
	return &ServiceImpl{
		client:             client,
		userDailyUsageRepo: userDailyUsageRepo,
		clock:              clock,
		limits:             limits,
	}
}

var (
	ErrQuizNotFound    = errors.New("word not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrDeleteQuiz      = errors.New("failed to delete word")
	ErrDatabaseFailure = errors.New("database failure")
)
