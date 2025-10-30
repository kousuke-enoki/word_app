package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"word_app/backend/src/domain"
	linemock "word_app/backend/src/mocks/infrastructure/auth/line"
	jwtmock "word_app/backend/src/mocks/infrastructure/jwt"
	authmock "word_app/backend/src/mocks/infrastructure/repository/auth"
	settingmock "word_app/backend/src/mocks/infrastructure/repository/setting"
	txmock "word_app/backend/src/mocks/infrastructure/repository/tx"
	usermock "word_app/backend/src/mocks/infrastructure/repository/user"
	udumock "word_app/backend/src/mocks/infrastructure/repository/userdailyusage"
	"word_app/backend/src/usecase/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClock struct {
	now time.Time
}

func (m *mockClock) Now() time.Time {
	return m.now
}

func makeAuthUC(t *testing.T, tm *txmock.MockManager, provider *linemock.MockProvider, userRepo *usermock.MockRepository, settingRepo *settingmock.MockUserConfigRepository, extAuthRepo *authmock.MockExternalAuthRepository, jwtGen *jwtmock.MockJWTGenerator, tempJwtGen *jwtmock.MockTempTokenGenerator, rootSettingRepo *settingmock.MockRootConfigRepository, userDailyUsageRepo *udumock.MockRepository, clock *mockClock) *auth.AuthUsecase {
	return auth.NewUsecase(tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock)
}

func TestAuthUsecase_TestLogin(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	clock := &mockClock{now: now}

	t.Run("success - test user created", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		// UUIDベースのユーザー名とメールアドレスが生成されるため、mock.Anythingを使用
		userRepo.On("Create", ctx, mock.Anything).
			Return(&domain.User{ID: 123, Name: "テストユーザー@abc12345", IsTest: true}, nil)
		settingRepo.On("CreateDefault", ctx, 123).Return(nil)
		userDailyUsageRepo.On("CreateIfNotExists", ctx, 123, now).Return(nil)

		result, err := uc.TestLogin(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 123, result.UserID)
		assert.Contains(t, result.UserName, "テストユーザー@")
		tm.AssertExpectations(t)
	})

	t.Run("error - test user mode disabled", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: false}, nil)

		result, err := uc.TestLogin(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test user mode is disabled")
		assert.Nil(t, result)
		tm.AssertNotCalled(t, "Begin", ctx)
	})

	t.Run("error - Get root setting fails", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock)

		rootSettingRepo.On("Get", ctx).
			Return(nil, errors.New("database error"))

		result, err := uc.TestLogin(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.Nil(t, result)
		tm.AssertNotCalled(t, "Begin", ctx)
	})

	t.Run("error - Begin fails", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)
		tm.On("Begin", ctx).Return(nil, nil, errors.New("tx error"))

		result, err := uc.TestLogin(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tx error")
		assert.Nil(t, result)
	})

	t.Run("error - Create user fails", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("Create", ctx, mock.Anything).
			Return(nil, errors.New("user create error"))

		result, err := uc.TestLogin(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user create error")
		assert.Nil(t, result)
	})

	t.Run("error - CreateDefault fails", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("Create", ctx, mock.Anything).
			Return(&domain.User{ID: 123, Name: "テストユーザー@abc12345"}, nil)
		settingRepo.On("CreateDefault", ctx, 123).Return(errors.New("setting error"))

		result, err := uc.TestLogin(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "setting error")
		assert.Nil(t, result)
	})

	t.Run("error - CreateIfNotExists fails", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("Create", ctx, mock.Anything).
			Return(&domain.User{ID: 123, Name: "テストユーザー@abc12345"}, nil)
		settingRepo.On("CreateDefault", ctx, 123).Return(nil)
		userDailyUsageRepo.On("CreateIfNotExists", ctx, 123, now).Return(errors.New("daily usage error"))

		result, err := uc.TestLogin(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "daily usage error")
		assert.Nil(t, result)
	})
}
