// backend/src/usecase/user/test/get_detail_with_mocks_test.go
package user_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"word_app/backend/src/domain"
	authmock "word_app/backend/src/mocks/infrastructure/repository/auth"
	settingmock "word_app/backend/src/mocks/infrastructure/repository/setting"
	txmock "word_app/backend/src/mocks/infrastructure/repository/tx"
	usermock "word_app/backend/src/mocks/infrastructure/repository/user"
	"word_app/backend/src/models"
	uc "word_app/backend/src/usecase/user"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T { return &v }

func fixedTime() time.Time {
	// 2025-01-02 03:04:05 +0000（UTC）に固定
	return time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
}

func mustFmt(t time.Time) string { return t.Format("2006-01-02 15:04:05") }
func makeUC(t *testing.T, ur *usermock.MockRepository) *uc.UserUsecase {
	ctx := context.Background()

	// tx.Manager の mock
	tm := txmock.NewMockManager(t)
	// Delete や本件の Get* 系は Begin を使わないなら期待値不要。
	// もし Begin を使うユースケースであれば、こう設定:
	tm.EXPECT().
		Begin(mock.Anything).
		Return(ctx, func(bool) error { return nil }, nil).
		Maybe() // 呼ばれない可能性もあるなら Maybe

	return uc.NewUserUsecase(
		tm,
		ur,
		settingmock.NewMockUserConfigRepository(t),
		authmock.NewMockExternalAuthRepository(t),
	)
}

// ---- GetMyDetail ------------------------------------------------------------

func TestUserUsecase_GetMyDetail_WithMocks(t *testing.T) {
	t.Run("OK: returns own detail mapped to DTO", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		ucase := makeUC(t, ur)

		now := fixedTime()
		email := "me@example.com"
		ent := &domain.User{
			ID:          10,
			Name:        "Alice",
			Email:       &email,
			IsAdmin:     true,
			IsRoot:      false,
			IsTest:      true,
			HasLine:     true,
			HasPassword: true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		ur.EXPECT().
			FindDetailByID(mock.Anything, 10).
			Return(ent, nil).
			Once()

		got, err := ucase.GetMyDetail(context.Background(), 10)
		require.NoError(t, err)

		require.Equal(t, &models.UserDetail{
			ID:               10,
			Name:             "Alice",
			Email:            &email,
			IsAdmin:          true,
			IsRoot:           false,
			IsTest:           true,
			IsLine:           true,
			IsSettedPassword: true,
			CreatedAt:        mustFmt(now),
			UpdatedAt:        mustFmt(now),
		}, got)
	})

	t.Run("OK: email nil, flags false, time format", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		ucase := makeUC(t, ur)

		now := fixedTime()
		ent := &domain.User{
			ID:          1,
			Name:        "",
			Email:       nil,
			IsAdmin:     false,
			IsRoot:      false,
			IsTest:      false,
			HasLine:     false,
			HasPassword: false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		ur.EXPECT().
			FindDetailByID(mock.Anything, 1).
			Return(ent, nil).
			Once()

		got, err := ucase.GetMyDetail(context.Background(), 1)
		require.NoError(t, err)

		require.Nil(t, got.Email)
		require.Equal(t, mustFmt(now), got.CreatedAt)
		require.Equal(t, mustFmt(now), got.UpdatedAt)
		require.False(t, got.IsAdmin)
		require.False(t, got.IsRoot)
		require.False(t, got.IsTest)
		require.False(t, got.IsLine)
		require.False(t, got.IsSettedPassword)
	})

	t.Run("NG: repo returns not found", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		ucase := makeUC(t, ur)

		ur.EXPECT().
			FindDetailByID(mock.Anything, 99).
			Return(nil, errors.New("user not found")).
			Once()

		got, err := ucase.GetMyDetail(context.Background(), 99)
		require.Error(t, err)
		require.Nil(t, got)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("NG: repo returns internal error", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		ucase := makeUC(t, ur)

		ur.EXPECT().
			FindDetailByID(mock.Anything, 10).
			Return(nil, errors.New("internal")).
			Once()

		got, err := ucase.GetMyDetail(context.Background(), 10)
		require.Error(t, err)
		require.Nil(t, got)
		require.Contains(t, err.Error(), "internal")
	})
}

// ---- GetDetailByID ----------------------------------------------------------

func TestUserUsecase_GetDetailByID_WithMocks(t *testing.T) {
	t.Run("OK: admin viewer can get other user's detail", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		ucase := makeUC(t, ur)

		now := fixedTime()
		ur.EXPECT().
			FindByID(mock.Anything, 1).
			Return(&domain.User{ID: 1, IsAdmin: true}, nil).
			Once()

		email := "bob@example.com"
		target := &domain.User{
			ID:          2,
			Name:        "Bob",
			Email:       &email,
			IsAdmin:     false,
			IsRoot:      false,
			IsTest:      true,
			HasLine:     true,
			HasPassword: false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		ur.EXPECT().
			FindDetailByID(mock.Anything, 2).
			Return(target, nil).
			Once()

		got, err := ucase.GetDetailByID(context.Background(), 1, 2)
		require.NoError(t, err)
		require.Equal(t, &models.UserDetail{
			ID:               2,
			Name:             "Bob",
			Email:            &email,
			IsAdmin:          false,
			IsRoot:           false,
			IsTest:           true,
			IsLine:           true,
			IsSettedPassword: false,
			CreatedAt:        mustFmt(now),
			UpdatedAt:        mustFmt(now),
		}, got)
	})

	t.Run("NG: viewer not found", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		ucase := makeUC(t, ur)

		ur.EXPECT().
			FindByID(mock.Anything, 10).
			Return(nil, errors.New("viewer not found")).
			Once()

		got, err := ucase.GetDetailByID(context.Background(), 10, 20)
		require.Error(t, err)
		require.Nil(t, got)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("NG: viewer not admin -> forbidden", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		ucase := makeUC(t, ur)

		ur.EXPECT().
			FindByID(mock.Anything, 10).
			Return(&domain.User{ID: 10, IsAdmin: false}, nil).
			Once()
		// 非管理者のため FindDetailByID は呼ばれない

		got, err := ucase.GetDetailByID(context.Background(), 10, 20)
		require.Error(t, err)
		require.Nil(t, got)
		require.Contains(t, strings.ToLower(err.Error()), "forbidden")
	})

	t.Run("NG: target not found", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		ucase := makeUC(t, ur)

		ur.EXPECT().
			FindByID(mock.Anything, 1).
			Return(&domain.User{ID: 1, IsAdmin: true}, nil).
			Once()
		ur.EXPECT().
			FindDetailByID(mock.Anything, 99).
			Return(nil, errors.New("user not found")).
			Once()

		got, err := ucase.GetDetailByID(context.Background(), 1, 99)
		require.Error(t, err)
		require.Nil(t, got)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("NG: internal error on FindDetailByID", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		ucase := makeUC(t, ur)

		ur.EXPECT().
			FindByID(mock.Anything, 1).
			Return(&domain.User{ID: 1, IsAdmin: true}, nil).
			Once()
		ur.EXPECT().
			FindDetailByID(mock.Anything, 2).
			Return(nil, errors.New("internal")).
			Once()

		got, err := ucase.GetDetailByID(context.Background(), 1, 2)
		require.Error(t, err)
		require.Nil(t, got)
		require.Contains(t, err.Error(), "internal")
	})
}
