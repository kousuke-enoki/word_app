// backend/src/usecase/user/test/list_users_with_mocks_test.go
package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"word_app/backend/src/domain"
	"word_app/backend/src/domain/repository"
	authmock "word_app/backend/src/mocks/infrastructure/repository/auth"
	settingmock "word_app/backend/src/mocks/infrastructure/repository/setting"
	txmock "word_app/backend/src/mocks/infrastructure/repository/tx"
	usermock "word_app/backend/src/mocks/infrastructure/repository/user"
	"word_app/backend/src/models"
	uc "word_app/backend/src/usecase/user"
)

// func fixedTime() time.Time { return time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC) }
// func mustFmt(t time.Time) string {
// 	return t.Format("2006-01-02 15:04:05")
// }

func makeUC_List(t *testing.T, ur *usermock.MockRepository) *uc.UserUsecase {
	// tx.Manager を使わないユースケースだが、New には渡す必要がある
	tm := txmock.NewMockManager(t)
	tm.EXPECT().Begin(mock.Anything).Maybe().Return(context.Background(), func(bool) error { return nil }, nil)
	tm.EXPECT().WithTx(mock.Anything, mock.AnythingOfType("func(context.Context) error")).Maybe().Return(nil)

	return uc.NewUserUsecase(
		tm,
		ur,
		settingmock.NewMockUserConfigRepository(t),
		authmock.NewMockExternalAuthRepository(t),
	)
}

func TestUserUsecase_ListUsers_WithMocks(t *testing.T) {
	ctx := context.Background()

	t.Run("OK: root viewer, paging/sort/search are passed and DTO mapped; total pages rounded up", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		u := makeUC_List(t, ur)

		// viewer は root
		ur.EXPECT().
			FindByID(mock.Anything, 100).
			Return(&domain.User{ID: 100, IsRoot: true}, nil).
			Once()

		// 入力
		in := uc.ListUsersInput{
			ViewerID: 100,
			Search:   "ali",
			SortBy:   "name",
			Order:    "desc",
			Page:     2,
			Limit:    3,
		}
		// 計算: offset=(2-1)*3=3
		offset := 3
		now := fixedTime()
		email1 := "alice@example.com"
		email2 := "bob@example.com"

		entUsers := []*domain.User{
			{
				ID:          1,
				Name:        "Alice",
				Email:       &email1,
				IsAdmin:     true,
				IsRoot:      false,
				IsTest:      true,
				HasLine:     true,
				HasPassword: true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          2,
				Name:        "Bob",
				Email:       &email2,
				IsAdmin:     false,
				IsRoot:      false,
				IsTest:      false,
				HasLine:     false,
				HasPassword: false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		// ListUsers に渡るフィルタの中身を厳密チェック
		expFilter := repository.UserListFilter{
			Search: "ali",
			SortBy: "name",
			Order:  "desc",
			Offset: offset,
			Limit:  3,
		}
		ur.EXPECT().
			ListUsers(mock.Anything, mock.MatchedBy(func(f repository.UserListFilter) bool {
				require.Equal(t, expFilter, f)
				return true
			})).
			Return(&repository.UserListResult{
				Users:      entUsers,
				TotalCount: 8, // totalPages = ceil(8/3) = 3
			}, nil).
			Once()

		got, err := u.ListUsers(ctx, in)
		require.NoError(t, err)

		require.Equal(t, 3, got.TotalPages)
		require.Equal(t, []models.User{
			{
				ID:               1,
				Name:             "Alice",
				IsAdmin:          true,
				IsRoot:           false,
				IsTest:           true,
				Email:            &email1,
				IsSettedPassword: true,
				IsLine:           true,
				CreatedAt:        mustFmt(now),
				UpdatedAt:        mustFmt(now),
			},
			{
				ID:               2,
				Name:             "Bob",
				IsAdmin:          false,
				IsRoot:           false,
				IsTest:           false,
				Email:            &email2,
				IsSettedPassword: false,
				IsLine:           false,
				CreatedAt:        mustFmt(now),
				UpdatedAt:        mustFmt(now),
			},
		}, got.Users)
	})

	t.Run("OK: defaults when page<=0 and limit<=0; email nil and flags false; totalPages=0 when totalCount=0", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		u := makeUC_List(t, ur)

		ur.EXPECT().
			FindByID(mock.Anything, 1).
			Return(&domain.User{ID: 1, IsRoot: true}, nil).
			Once()

		// page/limit は 0 → デフォルト page=1, limit=20, offset=0
		in := uc.ListUsersInput{
			ViewerID: 1,
			Page:     0,
			Limit:    0,
			Search:   "",
			SortBy:   "",
			Order:    "",
		}
		now := fixedTime()
		entUsers := []*domain.User{
			{
				ID:          10,
				Name:        "",
				Email:       nil,
				IsAdmin:     false,
				IsRoot:      false,
				IsTest:      false,
				HasLine:     false,
				HasPassword: false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}
		ur.EXPECT().
			ListUsers(mock.Anything, mock.MatchedBy(func(f repository.UserListFilter) bool {
				require.Equal(t, repository.UserListFilter{
					Search: "",
					SortBy: "",
					Order:  "",
					Offset: 0,
					Limit:  20,
				}, f)
				return true
			})).
			Return(&repository.UserListResult{
				Users:      entUsers,
				TotalCount: 0, // totalPages = 0
			}, nil).
			Once()

		got, err := u.ListUsers(ctx, in)
		require.NoError(t, err)

		require.Equal(t, 0, got.TotalPages)
		require.Len(t, got.Users, 1)
		require.Nil(t, got.Users[0].Email)
		require.False(t, got.Users[0].IsAdmin)
		require.False(t, got.Users[0].IsRoot)
		require.False(t, got.Users[0].IsTest)
		require.False(t, got.Users[0].IsLine)
		require.False(t, got.Users[0].IsSettedPassword)
		require.Equal(t, mustFmt(now), got.Users[0].CreatedAt)
		require.Equal(t, mustFmt(now), got.Users[0].UpdatedAt)
	})

	t.Run("NG: viewer not found (FindByID returns error)", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		u := makeUC_List(t, ur)

		ur.EXPECT().
			FindByID(mock.Anything, 9).
			Return(nil, errors.New("not found")).
			Once()

		got, err := u.ListUsers(ctx, uc.ListUsersInput{
			ViewerID: 9,
			Page:     1,
			Limit:    10,
		})
		require.Error(t, err)
		require.Nil(t, got)
		require.Contains(t, err.Error(), "not found")
		// ListUsers は呼ばれない
	})

	t.Run("NG: viewer is not root -> forbidden", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		u := makeUC_List(t, ur)

		ur.EXPECT().
			FindByID(mock.Anything, 2).
			Return(&domain.User{ID: 2, IsRoot: false}, nil).
			Once()

		got, err := u.ListUsers(ctx, uc.ListUsersInput{
			ViewerID: 2,
			Page:     1,
			Limit:    10,
		})
		require.Error(t, err)
		require.Nil(t, got)
		require.Contains(t, err.Error(), "forbidden")
	})

	t.Run("NG: repo.ListUsers returns internal error", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		u := makeUC_List(t, ur)

		ur.EXPECT().
			FindByID(mock.Anything, 3).
			Return(&domain.User{ID: 3, IsRoot: true}, nil).
			Once()

		ur.EXPECT().
			ListUsers(mock.Anything, mock.MatchedBy(func(f repository.UserListFilter) bool {
				// page/limit が <=0 の場合のデフォルトもカバーしておく
				require.Equal(t, repository.UserListFilter{
					Search: "q",
					SortBy: "email",
					Order:  "asc",
					Offset: 0,  // page=0 → page=1
					Limit:  20, // limit=0 → 20
				}, f)
				return true
			})).
			Return(nil, errors.New("internal")).
			Once()

		got, err := u.ListUsers(ctx, uc.ListUsersInput{
			ViewerID: 3,
			Search:   "q",
			SortBy:   "email",
			Order:    "asc",
			Page:     0,
			Limit:    0,
		})
		require.Error(t, err)
		require.Nil(t, got)
		require.Contains(t, err.Error(), "internal")
	})

	t.Run("OK: empty list but positive totalCount -> totalPages computed", func(t *testing.T) {
		ur := usermock.NewMockRepository(t)
		u := makeUC_List(t, ur)

		ur.EXPECT().
			FindByID(mock.Anything, 50).
			Return(&domain.User{ID: 50, IsRoot: true}, nil).
			Once()

		// page=3, limit=5 → offset=10
		in := uc.ListUsersInput{
			ViewerID: 50, Page: 3, Limit: 5,
		}
		ur.EXPECT().
			ListUsers(mock.Anything, mock.MatchedBy(func(f repository.UserListFilter) bool {
				require.Equal(t, repository.UserListFilter{
					Search: "",
					SortBy: "",
					Order:  "",
					Offset: 10,
					Limit:  5,
				}, f)
				return true
			})).
			Return(&repository.UserListResult{
				Users:      []*domain.User{}, // 空
				TotalCount: 11,               // ceil(11/5)=3
			}, nil).
			Once()

		got, err := u.ListUsers(ctx, in)
		require.NoError(t, err)
		require.Equal(t, 3, got.TotalPages)
		require.Empty(t, got.Users)
	})
}
