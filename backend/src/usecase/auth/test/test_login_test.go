package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"word_app/backend/src/domain"
	"word_app/backend/src/infrastructure/ratelimit"
	linemock "word_app/backend/src/mocks/infrastructure/auth/line"
	jwtmock "word_app/backend/src/mocks/infrastructure/jwt"
	ratelimitmock "word_app/backend/src/mocks/infrastructure/ratelimit"
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

func makeAuthUC(t *testing.T, tm *txmock.MockManager, provider *linemock.MockProvider, userRepo *usermock.MockRepository, settingRepo *settingmock.MockUserConfigRepository, extAuthRepo *authmock.MockExternalAuthRepository, jwtGen *jwtmock.MockJWTGenerator, tempJwtGen *jwtmock.MockTempTokenGenerator, rootSettingRepo *settingmock.MockRootConfigRepository, userDailyUsageRepo *udumock.MockRepository, clock *mockClock, rateLimiter *ratelimitmock.MockRateLimiter) *auth.AuthUsecase {
	return auth.NewUsecase(tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)
}

func TestAuthUsecase_TestLoginWithRateLimit(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	clock := &mockClock{now: now}

	testIP := "192.168.1.1"
	testUAHash := "test-ua-hash"
	testRoute := "/users/auth/test-login"
	testJump := "quiz"

	t.Run("success - test user created (rate limit passed, no cache)", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		// モック設定
		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)

		// レート制限パス（キャッシュなし）
		rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
			Return(&ratelimit.RateLimitResult{
				Allowed:      true,
				LastPayload:  nil,
				CurrentCount: 1,
			}, nil)

		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("Create", ctx, mock.Anything).
			Return(&domain.User{ID: 123, Name: "テストユーザー@abc12345", IsTest: true}, nil)
		settingRepo.On("CreateDefault", ctx, 123).Return(nil)
		userDailyUsageRepo.On("CreateIfNotExists", ctx, 123, now).Return(nil)

		// JWT生成
		jwtGen.On("GenerateJWT", "123").Return("test-jwt-token", nil)

		// 成功レスポンスの保存
		rateLimiter.On("SaveLastResult", ctx, testIP, testUAHash, testRoute, mock.Anything).
			Return(nil)

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Nil(t, lastPayload)
		assert.Equal(t, 0, retryAfter)
		assert.Equal(t, 123, result.UserID)
		assert.Equal(t, "test-jwt-token", result.Token)
		assert.Contains(t, result.UserName, "テストユーザー@")
		assert.Equal(t, "quiz", result.Jump)

		tm.AssertExpectations(t)
		rateLimiter.AssertExpectations(t)
		jwtGen.AssertExpectations(t)
	})

	t.Run("success - cached response returned (rate limit passed, cache exists)", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)

		// レート制限パス + キャッシュあり
		cachedPayload := []byte(`{"token":"cached-token","user_id":999,"user_name":"Cached User","jump":"list"}`)
		rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
			Return(&ratelimit.RateLimitResult{
				Allowed:      true,
				LastPayload:  cachedPayload,
				CurrentCount: 3,
			}, nil)

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.NoError(t, err)
		assert.Nil(t, result)
		assert.Equal(t, cachedPayload, lastPayload)
		assert.Equal(t, 0, retryAfter)

		// DB操作は行われないことを確認
		tm.AssertNotCalled(t, "Begin")
		userRepo.AssertNotCalled(t, "Create")
		rateLimiter.AssertExpectations(t)
	})

	t.Run("success - cached response returned (rate limit exceeded, cache exists)", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)

		// レート制限超過 + キャッシュあり（Allowed: trueだが実質キャッシュ返却）
		cachedPayload := []byte(`{"token":"cached-token","user_id":888,"user_name":"Rate Limited User","jump":"bulk"}`)
		rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
			Return(&ratelimit.RateLimitResult{
				Allowed:      true, // DynamoDB実装ではキャッシュがあればAllowed=true
				LastPayload:  cachedPayload,
				CurrentCount: 10,
			}, nil)

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.NoError(t, err)
		assert.Nil(t, result)
		assert.Equal(t, cachedPayload, lastPayload)
		assert.Equal(t, 0, retryAfter)

		tm.AssertNotCalled(t, "Begin")
		rateLimiter.AssertExpectations(t)
	})

	t.Run("error - rate limit exceeded (no cache)", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)

		// レート制限超過 + キャッシュなし
		rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
			Return(&ratelimit.RateLimitResult{
				Allowed:      false,
				LastPayload:  nil,
				RetryAfter:   45,
				CurrentCount: 10,
			}, nil)

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate limited")
		assert.Nil(t, result)
		assert.Nil(t, lastPayload)
		assert.Equal(t, 45, retryAfter)

		tm.AssertNotCalled(t, "Begin")
		rateLimiter.AssertExpectations(t)
	})

	t.Run("error - CheckRateLimit fails", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)

		rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
			Return(nil, errors.New("dynamodb connection error"))

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dynamodb connection error")
		assert.Nil(t, result)
		assert.Nil(t, lastPayload)
		assert.Equal(t, 0, retryAfter)

		tm.AssertNotCalled(t, "Begin")
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
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: false}, nil)

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test user mode is disabled")
		assert.Nil(t, result)
		assert.Nil(t, lastPayload)
		assert.Equal(t, 0, retryAfter)

		tm.AssertNotCalled(t, "Begin")
		rateLimiter.AssertNotCalled(t, "CheckRateLimit")
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
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(nil, errors.New("database error"))

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.Nil(t, result)
		assert.Nil(t, lastPayload)
		assert.Equal(t, 0, retryAfter)

		tm.AssertNotCalled(t, "Begin")
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
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)
		rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
			Return(&ratelimit.RateLimitResult{Allowed: true}, nil)
		tm.On("Begin", ctx).Return(nil, nil, errors.New("tx error"))

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tx error")
		assert.Nil(t, result)
		assert.Nil(t, lastPayload)
		assert.Equal(t, 0, retryAfter)
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
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)
		rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
			Return(&ratelimit.RateLimitResult{Allowed: true}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("Create", ctx, mock.Anything).
			Return(nil, errors.New("user create error"))

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user create error")
		assert.Nil(t, result)
		assert.Nil(t, lastPayload)
		assert.Equal(t, 0, retryAfter)
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
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)
		rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
			Return(&ratelimit.RateLimitResult{Allowed: true}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("Create", ctx, mock.Anything).
			Return(&domain.User{ID: 123, Name: "テストユーザー@abc12345"}, nil)
		settingRepo.On("CreateDefault", ctx, 123).Return(errors.New("setting error"))

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "setting error")
		assert.Nil(t, result)
		assert.Nil(t, lastPayload)
		assert.Equal(t, 0, retryAfter)
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
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)
		rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
			Return(&ratelimit.RateLimitResult{Allowed: true}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("Create", ctx, mock.Anything).
			Return(&domain.User{ID: 123, Name: "テストユーザー@abc12345"}, nil)
		settingRepo.On("CreateDefault", ctx, 123).Return(nil)
		userDailyUsageRepo.On("CreateIfNotExists", ctx, 123, now).Return(errors.New("daily usage error"))

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "daily usage error")
		assert.Nil(t, result)
		assert.Nil(t, lastPayload)
		assert.Equal(t, 0, retryAfter)
	})

	t.Run("error - GenerateJWT fails", func(t *testing.T) {
		tm := txmock.NewMockManager(t)
		provider := linemock.NewMockProvider(t)
		userRepo := usermock.NewMockRepository(t)
		settingRepo := settingmock.NewMockUserConfigRepository(t)
		extAuthRepo := authmock.NewMockExternalAuthRepository(t)
		jwtGen := jwtmock.NewMockJWTGenerator(t)
		tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
		rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)
		rateLimiter := ratelimitmock.NewMockRateLimiter(t)

		uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

		rootSettingRepo.On("Get", ctx).
			Return(&domain.RootConfig{IsTestUserMode: true}, nil)
		rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
			Return(&ratelimit.RateLimitResult{Allowed: true}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("Create", ctx, mock.Anything).
			Return(&domain.User{ID: 123, Name: "テストユーザー@abc12345"}, nil)
		settingRepo.On("CreateDefault", ctx, 123).Return(nil)
		userDailyUsageRepo.On("CreateIfNotExists", ctx, 123, now).Return(nil)
		jwtGen.On("GenerateJWT", "123").Return("", errors.New("jwt generation error"))

		result, lastPayload, retryAfter, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, testJump)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "jwt generation error")
		assert.Nil(t, result)
		assert.Nil(t, lastPayload)
		assert.Equal(t, 0, retryAfter)
	})

	t.Run("success - jump parameter normalization", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"list is preserved", "list", "list"},
			{"bulk is preserved", "bulk", "bulk"},
			{"quiz is preserved", "quiz", "quiz"},
			{"invalid defaults to quiz", "invalid", "quiz"},
			{"empty defaults to quiz", "", "quiz"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tm := txmock.NewMockManager(t)
				provider := linemock.NewMockProvider(t)
				userRepo := usermock.NewMockRepository(t)
				settingRepo := settingmock.NewMockUserConfigRepository(t)
				extAuthRepo := authmock.NewMockExternalAuthRepository(t)
				jwtGen := jwtmock.NewMockJWTGenerator(t)
				tempJwtGen := jwtmock.NewMockTempTokenGenerator(t)
				rootSettingRepo := settingmock.NewMockRootConfigRepository(t)
				userDailyUsageRepo := udumock.NewMockRepository(t)
				rateLimiter := ratelimitmock.NewMockRateLimiter(t)

				uc := makeAuthUC(t, tm, provider, userRepo, settingRepo, extAuthRepo, jwtGen, tempJwtGen, rootSettingRepo, userDailyUsageRepo, clock, rateLimiter)

				rootSettingRepo.On("Get", ctx).
					Return(&domain.RootConfig{IsTestUserMode: true}, nil)
				rateLimiter.On("CheckRateLimit", ctx, testIP, testUAHash, testRoute).
					Return(&ratelimit.RateLimitResult{Allowed: true}, nil)
				tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
				userRepo.On("Create", ctx, mock.Anything).
					Return(&domain.User{ID: 123, Name: "テストユーザー@abc12345"}, nil)
				settingRepo.On("CreateDefault", ctx, 123).Return(nil)
				userDailyUsageRepo.On("CreateIfNotExists", ctx, 123, now).Return(nil)
				jwtGen.On("GenerateJWT", "123").Return("test-token", nil)
				rateLimiter.On("SaveLastResult", ctx, testIP, testUAHash, testRoute, mock.Anything).
					Return(nil)

				result, _, _, err := uc.TestLoginWithRateLimit(ctx, testIP, testUAHash, testRoute, tt.input)

				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result.Jump)
			})
		}
	})
}
