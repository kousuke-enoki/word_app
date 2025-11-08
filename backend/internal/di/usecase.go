// internal/di/usecase.go
package di

import (
	"word_app/backend/config"
	"word_app/backend/src/infrastructure/auth/line"
	"word_app/backend/src/infrastructure/jwt"
	"word_app/backend/src/infrastructure/ratelimit"
	authUc "word_app/backend/src/usecase/auth"
	bulkUc "word_app/backend/src/usecase/bulk"
	"word_app/backend/src/usecase/clock"
	jwtUc "word_app/backend/src/usecase/jwt"
	settingUc "word_app/backend/src/usecase/setting"
	userUc "word_app/backend/src/usecase/user"
	"word_app/backend/src/utils/tempjwt"
)

type UseCases struct {
	Auth         *authUc.AuthUsecase
	BulkToken    bulkUc.TokenizeUsecase
	BulkRegister bulkUc.RegisterUsecase
	Setting      settingUc.SettingFacade // interface
	User         *userUc.UserUsecase     // interface
	Jwt          *jwtUc.JwtUsecase       // interface
}

func NewUseCases(config *config.Config, r *Repos) (*UseCases, error) {
	// -------- Auth -----------
	// LINE 認証プロバイダの初期化
	lineProv, _ := line.NewProvider(config.Line)
	// JWT 生成器の初期化
	// JWTSecret は環境変数から取得することを想定
	jwtGen := jwt.NewMyJWTGenerator(config.JWT.Secret)
	tempJwt := tempjwt.New(config.JWT.TempSecret)
	rl, err := ratelimit.NewRateLimiterFromEnv()
	if err != nil {
		return nil, err
	}
	// -------- Setting -----------
	// 各種設定ユースケースの初期化
	// authCfgUc := settingUc.NewAuthConfig(r.RootSetting)
	getRootUc := settingUc.NewGetRootConfig(r.User, r.RootSetting)
	getRuntimeConfigUc := settingUc.NewRuntimeConfig(r.RootSetting, clock.SystemClock{})
	getUserUc := settingUc.NewGetUserConfig(r.UserSetting)
	updateRootUc := settingUc.NewUpdateRootConfig(r.RootSetting, r.User)
	updateUserUc := settingUc.NewUpdateUserConfig(r.Tx, r.UserSetting)
	settingFacade := settingUc.NewSettingFacade(getRootUc, getRuntimeConfigUc, getUserUc, updateRootUc, updateUserUc)

	return &UseCases{
		Auth: authUc.NewUsecase(r.Tx, lineProv, r.User, r.UserSetting,
			r.Auth, jwtGen, tempJwt, r.RootSetting, r.UserDailyUsage, clock.SystemClock{}, rl),

		BulkToken:    bulkUc.NewTokenizeUsecase(r.WordRead, r.RegisteredWordRead, r.UserDailyUsage, clock.SystemClock{}, &config.Limits),
		BulkRegister: bulkUc.NewRegisterUsecase(r.WordRead, r.RegisteredWordRead, r.RegisteredWordWrite, r.Tx, r.User, &config.Limits),

		Setting: settingFacade, // まとめ役だけ保持
		User:    userUc.NewUserUsecase(r.Tx, r.User, r.UserSetting, r.Auth),
		Jwt:     jwtUc.NewJwtUsecase(jwt.NewHS256Verifier(config.JWT.Secret), r.User),
	}, nil
}
