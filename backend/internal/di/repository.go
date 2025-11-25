// internal/di/repository.go
package di

import (
	authRepo "word_app/backend/src/infrastructure/repository/auth"
	registeredwordRepo "word_app/backend/src/infrastructure/repository/registeredword"
	settingRepo "word_app/backend/src/infrastructure/repository/setting"
	txRepo "word_app/backend/src/infrastructure/repository/tx"
	userRepo "word_app/backend/src/infrastructure/repository/user"
	userdailyusageRepo "word_app/backend/src/infrastructure/repository/userdailyusage"
	wordRepo "word_app/backend/src/infrastructure/repository/word"
	"word_app/backend/src/interfaces"
	"word_app/backend/src/interfaces/sqlexec"
)

type Repos struct {
	Tx                  txRepo.Manager // 既存の Tx ラッパーを流用
	User                userRepo.Repository
	Auth                authRepo.ExternalAuthRepository
	RootSetting         settingRepo.RootConfigRepository
	UserSetting         settingRepo.UserConfigRepository
	UserDailyUsage      userdailyusageRepo.Repository
	WordRead            wordRepo.ReadRepository
	RegisteredWordRead  registeredwordRepo.ReadRepository
	RegisteredWordWrite registeredwordRepo.WriteRepository
}

func NewRepositories(cli interfaces.ClientInterface, s sqlexec.Runner) *Repos {
	return &Repos{
		Tx:                  txRepo.NewEntManager(cli),
		User:                userRepo.NewEntUserRepo(cli),
		Auth:                authRepo.NewEntExtAuthRepo(cli),
		RootSetting:         settingRepo.NewEntRootConfigRepo(cli),
		UserSetting:         settingRepo.NewEntUserConfigRepo(cli),
		UserDailyUsage:      userdailyusageRepo.NewEntUserDailyUsageRepo(cli, s),
		WordRead:            wordRepo.NewEntWordReadRepo(cli),
		RegisteredWordRead:  registeredwordRepo.NewEntRegisteredWordReadRepo(cli),
		RegisteredWordWrite: registeredwordRepo.NewEntRegisteredWordWriteRepo(cli),
	}
}
