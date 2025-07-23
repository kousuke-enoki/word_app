package settinguctest

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"word_app/backend/src/domain"
	settingUc "word_app/backend/src/usecase/setting"

	// mockery 生成パッケージを調整
	mockAuth "word_app/backend/src/mocks/usecase/setting"
	mockRoot "word_app/backend/src/mocks/usecase/setting"
	mockUpdR "word_app/backend/src/mocks/usecase/setting"
	mockUpdU "word_app/backend/src/mocks/usecase/setting"
	mockUser "word_app/backend/src/mocks/usecase/setting"
)

/* shared data -------------------------------------------------------------- */
var (
	ctx     = context.Background()
	errFoo  = errors.New("boom")
	rootCfg = &domain.RootConfig{ID: 1, EditingPermission: "admin"}
	userCfg = &domain.UserConfig{ID: 1, UserID: 99, IsDarkMode: true}
	rootOut = &settingUc.OutputGetRootConfig{Config: rootCfg}
	userOut = &settingUc.OutputGetUserConfig{Config: userCfg}
	authDTO = &settingUc.AuthConfigDTO{IsLineAuth: true}
)

/* -------------------------------------------------------------------------- */
/*                                GetAuth                                     */
/* -------------------------------------------------------------------------- */

func TestSettingFacade_GetAuth(t *testing.T) {
	// success
	{
		a := mockAuth.NewMockGetAuthConfig(t)
		a.On("Execute", ctx).Return(authDTO, nil)
		f := settingUc.NewSettingFacade(a,
			mockRoot.NewMockGetRootConfig(t),
			mockUser.NewMockGetUserConfig(t),
			mockUpdR.NewMockUpdateRootConfig(t),
			mockUpdU.NewMockUpdateUserConfig(t),
		)
		got, err := f.GetAuth(ctx)
		assert.NoError(t, err)
		assert.Equal(t, true, got.IsLineAuth)
		a.AssertExpectations(t)
	}

	// error
	{
		a := mockAuth.NewMockGetAuthConfig(t)
		a.On("Execute", ctx).Return((*settingUc.AuthConfigDTO)(nil), errFoo)

		f := settingUc.NewSettingFacade(a,
			mockRoot.NewMockGetRootConfig(t),
			mockUser.NewMockGetUserConfig(t),
			mockUpdR.NewMockUpdateRootConfig(t),
			mockUpdU.NewMockUpdateUserConfig(t),
		)
		_, err := f.GetAuth(ctx)
		assert.ErrorIs(t, err, errFoo)
		a.AssertExpectations(t)
	}
}

/* -------------------------------------------------------------------------- */
/*                                GetRoot                                     */
/* -------------------------------------------------------------------------- */

func TestSettingFacade_GetRoot(t *testing.T) {
	in := settingUc.InputGetRootConfig{UserID: 99}

	// success
	{
		gr := mockRoot.NewMockGetRootConfig(t)
		gr.On("Execute", ctx, in).Return(rootOut, nil)

		f := settingUc.NewSettingFacade(
			mockAuth.NewMockGetAuthConfig(t),
			gr,
			mockUser.NewMockGetUserConfig(t),
			mockUpdR.NewMockUpdateRootConfig(t),
			mockUpdU.NewMockUpdateUserConfig(t),
		)
		out, err := f.GetRoot(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, "admin", out.Config.EditingPermission)
		gr.AssertExpectations(t)
	}

	// error
	{
		gr := mockRoot.NewMockGetRootConfig(t)
		gr.On("Execute", ctx, in).Return((*settingUc.OutputGetRootConfig)(nil), errFoo)
		f := settingUc.NewSettingFacade(
			mockAuth.NewMockGetAuthConfig(t),
			gr,
			mockUser.NewMockGetUserConfig(t),
			mockUpdR.NewMockUpdateRootConfig(t),
			mockUpdU.NewMockUpdateUserConfig(t),
		)
		_, err := f.GetRoot(ctx, in)
		assert.ErrorIs(t, err, errFoo)
		gr.AssertExpectations(t)
	}
}

/* -------------------------------------------------------------------------- */
/*                                GetUser                                     */
/* -------------------------------------------------------------------------- */

func TestSettingFacade_GetUser(t *testing.T) {
	in := settingUc.InputGetUserConfig{UserID: 99}

	// success
	{
		gu := mockUser.NewMockGetUserConfig(t)
		gu.On("Execute", ctx, in).Return(userOut, nil)
		f := settingUc.NewSettingFacade(
			mockAuth.NewMockGetAuthConfig(t),
			mockRoot.NewMockGetRootConfig(t),
			gu,
			mockUpdR.NewMockUpdateRootConfig(t),
			mockUpdU.NewMockUpdateUserConfig(t),
		)
		out, err := f.GetUser(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, true, out.Config.IsDarkMode)
		gu.AssertExpectations(t)
	}

	// error
	{
		gu := mockUser.NewMockGetUserConfig(t)
		gu.On("Execute", ctx, in).Return((*settingUc.OutputGetUserConfig)(nil), errFoo)
		f := settingUc.NewSettingFacade(
			mockAuth.NewMockGetAuthConfig(t),
			mockRoot.NewMockGetRootConfig(t),
			gu,
			mockUpdR.NewMockUpdateRootConfig(t),
			mockUpdU.NewMockUpdateUserConfig(t),
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
		ur := mockUpdR.NewMockUpdateRootConfig(t)
		ur.On("Execute", ctx, in).Return(rootCfg, nil)
		f := settingUc.NewSettingFacade(
			mockAuth.NewMockGetAuthConfig(t),
			mockRoot.NewMockGetRootConfig(t),
			mockUser.NewMockGetUserConfig(t),
			ur,
			mockUpdU.NewMockUpdateUserConfig(t),
		)
		out, err := f.UpdateRoot(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, "admin", out.EditingPermission)
		ur.AssertExpectations(t)
	}

	// error
	{
		ur := mockUpdR.NewMockUpdateRootConfig(t)
		ur.On("Execute", ctx, in).Return((*domain.RootConfig)(nil), errFoo)
		f := settingUc.NewSettingFacade(
			mockAuth.NewMockGetAuthConfig(t),
			mockRoot.NewMockGetRootConfig(t),
			mockUser.NewMockGetUserConfig(t),
			ur,
			mockUpdU.NewMockUpdateUserConfig(t),
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
		uu := mockUpdU.NewMockUpdateUserConfig(t)
		uu.On("Execute", ctx, in).Return(userCfg, nil)
		f := settingUc.NewSettingFacade(
			mockAuth.NewMockGetAuthConfig(t),
			mockRoot.NewMockGetRootConfig(t),
			mockUser.NewMockGetUserConfig(t),
			mockUpdR.NewMockUpdateRootConfig(t),
			uu,
		)
		out, err := f.UpdateUser(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, true, out.IsDarkMode)
		uu.AssertExpectations(t)
	}

	// error
	{
		uu := mockUpdU.NewMockUpdateUserConfig(t)
		uu.On("Execute", ctx, in).Return((*domain.UserConfig)(nil), errFoo)
		f := settingUc.NewSettingFacade(
			mockAuth.NewMockGetAuthConfig(t),
			mockRoot.NewMockGetRootConfig(t),
			mockUser.NewMockGetUserConfig(t),
			mockUpdR.NewMockUpdateRootConfig(t),
			uu,
		)
		_, err := f.UpdateUser(ctx, in)
		assert.ErrorIs(t, err, errFoo)
		uu.AssertExpectations(t)
	}
}
