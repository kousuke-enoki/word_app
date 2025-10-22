package word

import (
	"word_app/backend/src/infrastructure/repository/userdailyusage"
	"word_app/backend/src/interfaces/http/word"
	"word_app/backend/src/usecase/clock"
)

type Handler struct {
	wordService        word.Service
	userDailyUsageRepo userdailyusage.Repository
	clock              clock.Clock
}

func NewHandler(
	wordService word.Service,
	userDailyUsageRepo userdailyusage.Repository,
	clock clock.Clock,
) *Handler {
	return &Handler{
		wordService:        wordService,
		userDailyUsageRepo: userDailyUsageRepo,
		clock:              clock,
	}
}
