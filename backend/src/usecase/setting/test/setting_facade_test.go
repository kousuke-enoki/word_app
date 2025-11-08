package settinguctest

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"word_app/backend/src/domain"
	settingUc "word_app/backend/src/usecase/setting"

	// mockery 生成パッケージを調整
	mockSetting "word_app/backend/src/mocks/usecase/setting"
)

/* shared data -------------------------------------------------------------- */
var (
	ctx     = context.Background()
	errFoo  = errors.New("boom")
	rootCfg = &domain.RootConfig{ID: 1, EditingPermission: "admin"}
	userCfg = &domain.UserConfig{ID: 1, UserID: 99, IsDarkMode: true}
	userOut = &settingUc.OutputGetUserConfig{Config: userCfg}
	authDTO = &settingUc.RuntimeConfigDTO{IsLineAuthentication: true}
)

/* -------------------------------------------------------------------------- */
/*                                GetRuntimeConfig                                     */
/* -------------------------------------------------------------------------- */

func TestSettingFacade_GetRuntimeConfig(t *testing.T) {
	// success
	{
		a := mockSetting.NewMockGetRuntimeConfig(t)
		a.On("Execute", ctx).Return(authDTO, nil)
		f := settingUc.NewSettingFacade(
			mockSetting.NewMockGetRootConfig(t),
			a,
			mockSetting.NewMockGetUserConfig(t),
			mockSetting.NewMockUpdateRootConfig(t),
			mockSetting.NewMockUpdateUserConfig(t),
		)
		got, err := f.GetRuntimeConfig(ctx)
		assert.NoError(t, err)
		assert.Equal(t, true, got.IsLineAuthentication)
		a.AssertExpectations(t)
	}

	// error
	{
		a := mockSetting.NewMockGetRuntimeConfig(t)
		a.On("Execute", ctx).Return((*settingUc.RuntimeConfigDTO)(nil), errFoo)

		f := settingUc.NewSettingFacade(
			mockSetting.NewMockGetRootConfig(t),
			a,
			mockSetting.NewMockGetUserConfig(t),
			mockSetting.NewMockUpdateRootConfig(t),
			mockSetting.NewMockUpdateUserConfig(t),
		)
		_, err := f.GetRuntimeConfig(ctx)
		assert.ErrorIs(t, err, errFoo)
		a.AssertExpectations(t)
	}
}

/* -------------------------------------------------------------------------- */
/*                                GetRoot                                     */
/* -------------------------------------------------------------------------- */

func TestSettingFacade_GetRoot(t *testing.T) {
	in := settingUc.InputGetUserConfig{UserID: 99}

	// success
	{
		gu := mockSetting.NewMockGetUserConfig(t)
		gu.On("Execute", ctx, in).Return(userOut, nil)
		f := settingUc.NewSettingFacade(
			mockSetting.NewMockGetRootConfig(t),
			mockSetting.NewMockGetRuntimeConfig(t),
			gu,
			mockSetting.NewMockUpdateRootConfig(t),
			mockSetting.NewMockUpdateUserConfig(t),
		)
		out, err := f.GetUser(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, true, out.Config.IsDarkMode)
		gu.AssertExpectations(t)
	}

	// error
	{
		gu := mockSetting.NewMockGetUserConfig(t)
		gu.On("Execute", ctx, in).Return((*settingUc.OutputGetUserConfig)(nil), errFoo)
		f := settingUc.NewSettingFacade(
			mockSetting.NewMockGetRootConfig(t),
			mockSetting.NewMockGetRuntimeConfig(t),
			gu,
			mockSetting.NewMockUpdateRootConfig(t),
			mockSetting.NewMockUpdateUserConfig(t),
		)
		_, err := f.GetUser(ctx, in)
		assert.ErrorIs(t, err, errFoo)
		gu.AssertExpectations(t)
	}
}

/* -------------------------------------------------------------------------- */
/*                               UpdateRoot                                   */
/* -------------------------------------------------------------------------- */

func TestSettingFacade_UpdateRoot(t *testing.T) {
	in := settingUc.InputUpdateRootConfig{UserID: 99}

	// success
	{
		ur := mockSetting.NewMockUpdateRootConfig(t)
		ur.On("Execute", ctx, in).Return(rootCfg, nil)
		f := settingUc.NewSettingFacade(
			mockSetting.NewMockGetRootConfig(t),
			mockSetting.NewMockGetRuntimeConfig(t),
			mockSetting.NewMockGetUserConfig(t),
			ur,
			mockSetting.NewMockUpdateUserConfig(t),
		)
		out, err := f.UpdateRoot(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, "admin", out.EditingPermission)
		ur.AssertExpectations(t)
	}

	// error
	{
		ur := mockSetting.NewMockUpdateRootConfig(t)
		ur.On("Execute", ctx, in).Return((*domain.RootConfig)(nil), errFoo)
		f := settingUc.NewSettingFacade(
			mockSetting.NewMockGetRootConfig(t),
			mockSetting.NewMockGetRuntimeConfig(t),
			mockSetting.NewMockGetUserConfig(t),
			ur,
			mockSetting.NewMockUpdateUserConfig(t),
		)
		_, err := f.UpdateRoot(ctx, in)
		assert.ErrorIs(t, err, errFoo)
		ur.AssertExpectations(t)
	}
}

/* -------------------------------------------------------------------------- */
/*                               UpdateUser                                   */
/* -------------------------------------------------------------------------- */

func TestSettingFacade_UpdateUser(t *testing.T) {
	in := settingUc.InputUpdateUserConfig{UserID: 99}

	// success
	{
		uu := mockSetting.NewMockUpdateUserConfig(t)
		uu.On("Execute", ctx, in).Return(userCfg, nil)
		f := settingUc.NewSettingFacade(
			mockSetting.NewMockGetRootConfig(t),
			mockSetting.NewMockGetRuntimeConfig(t),
			mockSetting.NewMockGetUserConfig(t),
			mockSetting.NewMockUpdateRootConfig(t),
			uu,
		)
		out, err := f.UpdateUser(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, true, out.IsDarkMode)
		uu.AssertExpectations(t)
	}

	// error
	{
		uu := mockSetting.NewMockUpdateUserConfig(t)
		uu.On("Execute", ctx, in).Return((*domain.UserConfig)(nil), errFoo)
		f := settingUc.NewSettingFacade(
			mockSetting.NewMockGetRootConfig(t),
			mockSetting.NewMockGetRuntimeConfig(t),
			mockSetting.NewMockGetUserConfig(t),
			mockSetting.NewMockUpdateRootConfig(t),
			uu,
		)
		_, err := f.UpdateUser(ctx, in)
		assert.ErrorIs(t, err, errFoo)
		uu.AssertExpectations(t)
	}
}
