// internal/di/repository.go
package di

import (
	authRepo "word_app/backend/src/infrastructure/repository/auth"
	settingRepo "word_app/backend/src/infrastructure/repository/setting"
	txRepo "word_app/backend/src/infrastructure/repository/tx"
	userRepo "word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/interfaces"
)

type Repos struct {
	Tx          txRepo.Manager // 既存の Tx ラッパーを流用
	User        userRepo.Repository
	Auth        authRepo.ExternalAuthRepository
	RootSetting settingRepo.RootConfigRepository
	UserSetting settingRepo.UserConfigRepository
}

func NewRepositories(cli interfaces.ClientInterface) *Repos {
	return &Repos{
		Tx:          txRepo.NewEntManager(cli),
		User:        userRepo.NewEntUserRepo(cli),
		Auth:        authRepo.NewEntExtAuthRepo(cli),
		RootSetting: settingRepo.NewEntRootConfigRepo(cli),
		UserSetting: settingRepo.NewEntUserConfigRepo(cli),
	}
}
