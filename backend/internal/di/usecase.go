// internal/di/usecase.go
package di

import (
	"word_app/backend/config"
	"word_app/backend/src/infrastructure/auth/line"
	"word_app/backend/src/infrastructure/jwt"
	authUc "word_app/backend/src/usecase/auth"
	settingUc "word_app/backend/src/usecase/setting"
	userUc "word_app/backend/src/usecase/user"
	"word_app/backend/src/utils/tempjwt"
)

type UseCases struct {
	Auth    *authUc.Usecase
	Setting settingUc.SettingFacade // interface
	User    *userUc.UserUsecase     // interface
}

func NewUseCases(config *config.Config, r *Repos) (*UseCases, error) {
	// -------- Auth -----------
	// LINE 認証プロバイダの初期化
	lineProv, _ := line.NewProvider(config.Line)
	// JWT 生成器の初期化
	// JWTSecret は環境変数から取得することを想定
	jwtGen := jwt.NewMyJWTGenerator(config.JWT.Secret)
	tempJwt := tempjwt.New(config.JWT.TempSecret)

	// -------- Setting -----------
	// 各種設定ユースケースの初期化
	authCfgUc := settingUc.NewAuthConfig(r.RootSetting)
	getRootUc := settingUc.NewGetRootConfig(r.User, r.RootSetting)
	getUserUc := settingUc.NewGetUserConfig(r.UserSetting)
	updateRootUc := settingUc.NewUpdateRootConfig(r.RootSetting, r.User)
	updateUserUc := settingUc.NewUpdateUserConfig(r.Tx, r.UserSetting)
	settingFacade := settingUc.NewSettingFacade(authCfgUc, getRootUc, getUserUc, updateRootUc, updateUserUc)

	return &UseCases{
		Auth:    authUc.NewUsecase(r.Tx, lineProv, r.User, r.UserSetting, r.Auth, jwtGen, tempJwt),
		Setting: settingFacade, // まとめ役だけ保持
		User:    userUc.NewUserUsecase(r.Tx, r.User, r.UserSetting, r.Auth),
	}, nil
}
