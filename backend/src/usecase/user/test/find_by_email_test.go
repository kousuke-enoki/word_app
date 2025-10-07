// backend/src/usecase/user/test/find_by_email_test.go
package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"word_app/backend/src/domain"
	repomock "word_app/backend/src/mocks/infrastructure/repository/user" // mockery 生成物
	ucpkg "word_app/backend/src/usecase/user"
)

// UC の薄いラッパ（本番の FindByEmail と同じロジック）
type ucThin struct {
	userRepo interface {
		FindActiveByEmail(ctx context.Context, email string) (*domain.User, error)
	}
}

func (uc *ucThin) FindByEmail(ctx context.Context, email string) (*ucpkg.FindByEmailOutput, error) {
	u, err := uc.userRepo.FindActiveByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &ucpkg.FindByEmailOutput{
		UserID:         u.ID,
		HashedPassword: u.Password,
		IsAdmin:        u.IsAdmin,
		IsRoot:         u.IsRoot,
		IsTest:         u.IsTest,
	}, nil
}

func TestUserUsecase_FindByEmail_WithMockery(t *testing.T) {
	type given struct {
		email string
		user  *domain.User
		err   error
	}
	type want struct {
		out       *ucpkg.FindByEmailOutput
		errSubStr string // 空なら成功想定
	}

	run := func(name string, g given, w want) {
		t.Run(name, func(t *testing.T) {
			// mockery 生成モック
			repo := repomock.NewMockRepository(t)
			// どの Context が来ても良い: mock.Anything で許容
			repo.EXPECT().
				FindActiveByEmail(mock.Anything, g.email).
				Return(g.user, g.err)

			uc := &ucThin{userRepo: repo}

			out, err := uc.FindByEmail(context.Background(), g.email)

			// 期待どおり呼ばれている（引数も一致）ことは、上記 EXPECT で検証される
			// NewMockRepository が t.Cleanup で AssertExpectations を自動実行してくれる

			if w.errSubStr == "" {
				require.NoError(t, err)
				require.NotNil(t, out)

				require.Equal(t, w.out.UserID, out.UserID)
				require.Equal(t, w.out.HashedPassword, out.HashedPassword)
				require.Equal(t, w.out.IsAdmin, out.IsAdmin)
				require.Equal(t, w.out.IsRoot, out.IsRoot)
				require.Equal(t, w.out.IsTest, out.IsTest)
			} else {
				require.Error(t, err)
				require.Nil(t, out)
				require.Contains(t, err.Error(), w.errSubStr)
			}
		})
	}

	// ---------- 成功: 典型 ----------
	run("OK_typical",
		given{
			email: "alice@example.com",
			user: &domain.User{
				ID:       10,
				Password: "hashed-pw",
				IsAdmin:  true,
				IsRoot:   false,
				IsTest:   true,
			},
		},
		want{
			out: &ucpkg.FindByEmailOutput{
				UserID:         10,
				HashedPassword: "hashed-pw",
				IsAdmin:        true,
				IsRoot:         false,
				IsTest:         true,
			},
		},
	)

	// ---------- 成功: 境界（空パスワード・全フラグ false） ----------
	run("OK_boundary_empty_password_and_flags_false",
		given{
			email: "bob@example.com",
			user: &domain.User{
				ID:       20,
				Password: "",
				IsAdmin:  false,
				IsRoot:   false,
				IsTest:   false,
			},
		},
		want{
			out: &ucpkg.FindByEmailOutput{
				UserID:         20,
				HashedPassword: "",
				IsAdmin:        false,
				IsRoot:         false,
				IsTest:         false,
			},
		},
	)

	// ---------- 失敗: NotFound 相当 ----------
	run("NG_not_found",
		given{
			email: "nobody@example.com",
			err:   errors.New("user not found"),
		},
		want{errSubStr: "user not found"},
	)

	// ---------- 失敗: Internal 相当 ----------
	run("NG_internal_error",
		given{
			email: "err@example.com",
			err:   errors.New("internal server error"),
		},
		want{errSubStr: "internal"},
	)
}
