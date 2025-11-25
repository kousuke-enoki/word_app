package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	linemock "word_app/backend/src/mocks/infrastructure/auth/line"
	jwtmock "word_app/backend/src/mocks/infrastructure/jwt"
	ratelimitmock "word_app/backend/src/mocks/infrastructure/ratelimit"
	authmock "word_app/backend/src/mocks/infrastructure/repository/auth"
	settingmock "word_app/backend/src/mocks/infrastructure/repository/setting"
	txmock "word_app/backend/src/mocks/infrastructure/repository/tx"
	usermock "word_app/backend/src/mocks/infrastructure/repository/user"
	udumock "word_app/backend/src/mocks/infrastructure/repository/userdailyusage"
	"word_app/backend/src/usecase/auth"
	"word_app/backend/src/usecase/clock"

	"github.com/stretchr/testify/assert"
)

func makeAuthUCForTestLogout(t *testing.T, tm *txmock.MockManager, provider *linemock.MockProvider, userRepo *usermock.MockRepository, settingRepo *settingmock.MockUserConfigRepository, extAuthRepo *authmock.MockExternalAuthRepository, jwtGen *jwtmock.MockJWTGenerator, tempJwtGen *jwtmock.MockTempTokenGenerator, rootSettingRepo *settingmock.MockRootConfigRepository, userDailyUsageRepo *udumock.MockRepository, clock clock.Clock, rateLimiter *ratelimitmock.MockRateLimiter) *auth.AuthUsecase {
	return auth.NewUsecase(tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)
}

func TestAuthUsecase_TestLogout(t *testing.T) {
	ctx := context.Background()

	t.Run("success - test user deleted", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		clock := &mockClock{now: time.Now()}
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUCForTestLogout(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("DeleteIfTest", ctx, 123).Return(true, nil)
		rateLimiter.On("ClearCacheForUser", ctx, 123).Return(nil)

		err := uc.TestLogout(ctx, 123)

		assert.NoError(t, err)
		tm.AssertExpectations(t)
	})

	t.Run("success - already deleted (idempotent)", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		clock := &mockClock{now: time.Now()}
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUCForTestLogout(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("DeleteIfTest", ctx, 999).Return(false, nil)
		userRepo.On("Exists", ctx, 999).Return(false, nil) // 既に削除済み

		err := uc.TestLogout(ctx, 999)

		assert.NoError(t, err) // 冪等性：既に削除済みでも成功
		tm.AssertExpectations(t)
	})

	t.Run("error - forbidden for non-test user", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		clock := &mockClock{now: time.Now()}
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUCForTestLogout(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("DeleteIfTest", ctx, 456).Return(false, nil)
		userRepo.On("Exists", ctx, 456).Return(true, nil)  // 存在する
		userRepo.On("IsTest", ctx, 456).Return(false, nil) // 非テストユーザー

		err := uc.TestLogout(ctx, 456)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only test user can be deleted via test-logout")
		tm.AssertExpectations(t)
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
		clock := &mockClock{now: time.Now()}
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUCForTestLogout(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		tm.On("Begin", ctx).Return(nil, nil, errors.New("tx error"))

		err := uc.TestLogout(ctx, 123)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tx error")
	})

	t.Run("error - DeleteIfTest fails", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		clock := &mockClock{now: time.Now()}
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUCForTestLogout(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("DeleteIfTest", ctx, 123).Return(false, errors.New("db error"))

		err := uc.TestLogout(ctx, 123)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("error - Exists fails", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		clock := &mockClock{now: time.Now()}
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUCForTestLogout(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("DeleteIfTest", ctx, 123).Return(false, nil)
		userRepo.On("Exists", ctx, 123).Return(false, errors.New("exists error"))

		err := uc.TestLogout(ctx, 123)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exists error")
	})

	t.Run("error - IsTest fails", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		clock := &mockClock{now: time.Now()}
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUCForTestLogout(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("DeleteIfTest", ctx, 789).Return(false, nil)
		userRepo.On("Exists", ctx, 789).Return(true, nil)
		userRepo.On("IsTest", ctx, 789).Return(false, errors.New("isTest error"))

		err := uc.TestLogout(ctx, 789)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "isTest error")
	})

	t.Run("success - test user exists but DeleteIfTest returns false", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		clock := &mockClock{now: time.Now()}
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUCForTestLogout(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("DeleteIfTest", ctx, 321).Return(false, nil)
		userRepo.On("Exists", ctx, 321).Return(true, nil)
		userRepo.On("IsTest", ctx, 321).Return(true, nil) // テストユーザーだが削除されなかった場合

		err := uc.TestLogout(ctx, 321)

		// 実装ではこの場合も成功扱い（commit = true）
		assert.NoError(t, err)
		tm.AssertExpectations(t)
	})
}
