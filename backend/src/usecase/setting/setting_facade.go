// src/usecase/setting/facade.go
package settingUc

import (
	"context"
	"word_app/backend/src/domain"
)

// Clean Architecture的に構造体は定義しつつ、使用時にはまとめて使用する
// Facade は設定関連のユースケースをまとめたinterface
// これを使うことで、設定関連の操作を一つのinterfaceで扱えるようにする
// 例えば、設定画面の API を実装する際に、
// ユーザー設定やルート設定を一つのinterfaceで扱えるようにするためのもの
// ただし、Facade はあくまでinterfaceであり、
// 実装は各ユースケースの実装を組み合わせて行う
type SettingFacade interface {
	GetAuth(ctx context.Context) (*AuthConfigDTO, error)
	GetRoot(ctx context.Context, in InputGetRootConfig) (*OutputGetRootConfig, error)
	GetUser(ctx context.Context, in InputGetUserConfig) (*OutputGetUserConfig, error)
	UpdateRoot(ctx context.Context, in InputUpdateRootConfig) (*domain.RootConfig, error)
	UpdateUser(ctx context.Context, in InputUpdateUserConfig) (*domain.UserConfig, error)
}

type settingFacade struct {
	authCfg    GetAuthConfig
	getRoot    GetRootConfig
	getUser    GetUserConfig
	updateRoot UpdateRootConfig
	updateUser UpdateUserConfig
}

func NewSettingFacade(
	a GetAuthConfig,
	gr GetRootConfig,
	gu GetUserConfig,
	ur UpdateRootConfig,
	uu UpdateUserConfig,
) SettingFacade {
	return &settingFacade{a, gr, gu, ur, uu}
}

// ↓ 各メソッドは単に委譲
func (f *settingFacade) GetAuth(ctx context.Context) (*AuthConfigDTO, error) {
	return f.authCfg.Execute(ctx)
}

func (f *settingFacade) GetRoot(ctx context.Context, in InputGetRootConfig) (*OutputGetRootConfig, error) {
	return f.getRoot.Execute(ctx, in)
}

func (f *settingFacade) GetUser(ctx context.Context, in InputGetUserConfig) (*OutputGetUserConfig, error) {
	return f.getUser.Execute(ctx, in)
}

func (f *settingFacade) UpdateRoot(ctx context.Context, in InputUpdateRootConfig) (*domain.RootConfig, error) {
	return f.updateRoot.Execute(ctx, in)
}

func (f *settingFacade) UpdateUser(ctx context.Context, in InputUpdateUserConfig) (*domain.UserConfig, error) {
	return f.updateUser.Execute(ctx, in)
}
