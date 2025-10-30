package jwt_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"word_app/backend/src/domain"
	jwtmock "word_app/backend/src/mocks/infrastructure/jwt"
	usermock "word_app/backend/src/mocks/infrastructure/repository/user"
	jwtusecase "word_app/backend/src/usecase/jwt"

	"github.com/stretchr/testify/assert"
)

func makeAuthUC(t *testing.T, verifier *jwtmock.MockTokenVerifier, userRepo *usermock.MockRepository) *jwtusecase.JwtUsecase {
	return jwtusecase.NewJwtUsecase(verifier, userRepo)
}

func TestJwtUsecase_Authenticate(t *testing.T) {
	ctx := context.Background()

	t.Run("success - normal user", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		verifier.On("VerifyAndExtractSubject", ctx, "valid_token").
			Return("123", nil)
		userRepo.On("FindByID", ctx, 123).
			Return(&domain.User{
				ID:        123,
				IsAdmin:   false,
				IsRoot:    false,
				IsTest:    false,
				DeletedAt: nil,
			}, nil)

		principal, err := uc.Authenticate(ctx, "valid_token")

		assert.NoError(t, err)
		assert.Equal(t, 123, principal.UserID)
		assert.False(t, principal.IsAdmin)
		assert.False(t, principal.IsRoot)
		assert.False(t, principal.IsTest)
		verifier.AssertExpectations(t)
	})

	t.Run("success - admin user", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		verifier.On("VerifyAndExtractSubject", ctx, "admin_token").
			Return("456", nil)
		userRepo.On("FindByID", ctx, 456).
			Return(&domain.User{
				ID:        456,
				IsAdmin:   true,
				IsRoot:    false,
				IsTest:    false,
				DeletedAt: nil,
			}, nil)

		principal, err := uc.Authenticate(ctx, "admin_token")

		assert.NoError(t, err)
		assert.Equal(t, 456, principal.UserID)
		assert.True(t, principal.IsAdmin)
		assert.False(t, principal.IsRoot)
		assert.False(t, principal.IsTest)
		verifier.AssertExpectations(t)
	})

	t.Run("success - root user", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		verifier.On("VerifyAndExtractSubject", ctx, "root_token").
			Return("789", nil)
		userRepo.On("FindByID", ctx, 789).
			Return(&domain.User{
				ID:        789,
				IsAdmin:   true,
				IsRoot:    true,
				IsTest:    false,
				DeletedAt: nil,
			}, nil)

		principal, err := uc.Authenticate(ctx, "root_token")

		assert.NoError(t, err)
		assert.Equal(t, 789, principal.UserID)
		assert.True(t, principal.IsAdmin)
		assert.True(t, principal.IsRoot)
		assert.False(t, principal.IsTest)
		verifier.AssertExpectations(t)
	})

	t.Run("success - test user", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		verifier.On("VerifyAndExtractSubject", ctx, "test_token").
			Return("999", nil)
		userRepo.On("FindByID", ctx, 999).
			Return(&domain.User{
				ID:        999,
				IsAdmin:   false,
				IsRoot:    false,
				IsTest:    true,
				DeletedAt: nil,
			}, nil)

		principal, err := uc.Authenticate(ctx, "test_token")

		assert.NoError(t, err)
		assert.Equal(t, 999, principal.UserID)
		assert.False(t, principal.IsAdmin)
		assert.False(t, principal.IsRoot)
		assert.True(t, principal.IsTest)
		verifier.AssertExpectations(t)
	})

	t.Run("error - VerifyAndExtractSubject fails", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		verifier.On("VerifyAndExtractSubject", ctx, "invalid_token").
			Return("", errors.New("token invalid"))

		principal, err := uc.Authenticate(ctx, "invalid_token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token invalid")
		assert.Equal(t, 0, principal.UserID)
		verifier.AssertExpectations(t)
		userRepo.AssertNotCalled(t, "FindByID", ctx, 123)
	})

	t.Run("error - invalid subject (not a number)", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		verifier.On("VerifyAndExtractSubject", ctx, "bad_subject_token").
			Return("not_a_number", nil)

		principal, err := uc.Authenticate(ctx, "bad_subject_token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unauthorized")
		assert.Equal(t, 0, principal.UserID)
		verifier.AssertExpectations(t)
		userRepo.AssertNotCalled(t, "FindByID", ctx, 123)
	})

	t.Run("error - user not found", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		verifier.On("VerifyAndExtractSubject", ctx, "notfound_token").
			Return("999", nil)
		userRepo.On("FindByID", ctx, 999).
			Return(nil, errors.New("user not found"))

		principal, err := uc.Authenticate(ctx, "notfound_token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.Equal(t, 0, principal.UserID)
		verifier.AssertExpectations(t)
	})

	t.Run("error - deleted user", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		now := time.Now()
		verifier.On("VerifyAndExtractSubject", ctx, "deleted_token").
			Return("111", nil)
		userRepo.On("FindByID", ctx, 111).
			Return(&domain.User{
				ID:        111,
				IsAdmin:   false,
				IsRoot:    false,
				IsTest:    false,
				DeletedAt: &now,
			}, nil)

		principal, err := uc.Authenticate(ctx, "deleted_token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unauthorized")
		assert.Equal(t, 0, principal.UserID)
		verifier.AssertExpectations(t)
	})

	t.Run("error - FindByID database error", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		verifier.On("VerifyAndExtractSubject", ctx, "db_error_token").
			Return("222", nil)
		userRepo.On("FindByID", ctx, 222).
			Return(nil, errors.New("database connection error"))

		principal, err := uc.Authenticate(ctx, "db_error_token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection error")
		assert.Equal(t, 0, principal.UserID)
		verifier.AssertExpectations(t)
	})

	t.Run("success - zero user ID", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		verifier.On("VerifyAndExtractSubject", ctx, "zero_token").
			Return("0", nil)
		userRepo.On("FindByID", ctx, 0).
			Return(&domain.User{
				ID:        0,
				IsAdmin:   false,
				IsRoot:    false,
				IsTest:    false,
				DeletedAt: nil,
			}, nil)

		principal, err := uc.Authenticate(ctx, "zero_token")

		assert.NoError(t, err)
		assert.Equal(t, 0, principal.UserID)
		assert.False(t, principal.IsAdmin)
		assert.False(t, principal.IsRoot)
		assert.False(t, principal.IsTest)
		verifier.AssertExpectations(t)
	})

	t.Run("error - empty subject", func(t *testing.T) {
		verifier := jwtmock.NewMockTokenVerifier(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeAuthUC(t, verifier, userRepo)

		verifier.On("VerifyAndExtractSubject", ctx, "empty_subject_token").
			Return("", nil)

		principal, err := uc.Authenticate(ctx, "empty_subject_token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unauthorized")
		assert.Equal(t, 0, principal.UserID)
		verifier.AssertExpectations(t)
		userRepo.AssertNotCalled(t, "FindByID", ctx, 123)
	})
}
