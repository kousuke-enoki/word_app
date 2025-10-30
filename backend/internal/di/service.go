// internal/di/handler.go
package di

import (
	"word_app/backend/config"
	"word_app/backend/src/interfaces"
	"word_app/backend/src/usecase/clock"

	quizSvc "word_app/backend/src/service/quiz"
	resultSvc "word_app/backend/src/service/result"
	wordSvc "word_app/backend/src/service/word"
)

type Services struct {
	Word   *wordSvc.ServiceImpl
	Quiz   *quizSvc.ServiceImpl
	Result *resultSvc.ServiceImpl
}

func NewServices(config *config.Config, uc *UseCases, client interfaces.ClientInterface, r *Repos) *Services {
	return &Services{
		Word: wordSvc.NewWordService(client, r.User,
			&config.Limits, r.Tx),
		Quiz:   quizSvc.NewService(client, r.UserDailyUsage, clock.SystemClock{}, &config.Limits),
		Result: resultSvc.NewService(client),
	}
}
