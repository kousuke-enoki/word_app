// src/usecase/user/test/update_user_test.go
package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"word_app/backend/src/domain"
	"word_app/backend/src/domain/repository"
	authmock "word_app/backend/src/mocks/infrastructure/repository/auth"
	settingmock "word_app/backend/src/mocks/infrastructure/repository/setting"
	txmock "word_app/backend/src/mocks/infrastructure/repository/tx"
	usermock "word_app/backend/src/mocks/infrastructure/repository/user"
	uc "word_app/backend/src/usecase/user"
)

func newUCWithMocks(t *testing.T, ur *usermock.MockRepository, tm *txmock.MockManager) *uc.UserUsecase {
	// setting/auth はこのユースケースでは呼ばれないのでダミー
	return uc.NewUserUsecase(
		tm,
		ur,
		settingmock.NewMockUserConfigRepository(t),
		authmock.NewMockExternalAuthRepository(t),
	)
}

func beginOK(ctx context.Context, tm *txmock.MockManager, onDone func(commit bool) error) {
	tm.EXPECT().
		Begin(mock.Anything).
		Return(ctx, onDone, nil).
		Once()
}

func beginErr(ctx context.Context, tm *txmock.MockManager, err error) {
	tm.EXPECT().
		Begin(mock.Anything).
		Return(ctx, nil, err).
		Once()
}

// func ptr[T any](v T) *T { return &v }

/************ ここから各ケースを独立関数に分割（中身は元のまま） ************/

func TestUpdateUser_SelfUpdateNameEmail_Success(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)

	editor := &domain.User{ID: 10, IsRoot: false, IsTest: false}
	target := &domain.User{ID: 10, IsRoot: false, IsTest: false}
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	ur.EXPECT().FindForUpdate(ctx, 10).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 10).Return(target, nil).Once()

	in := uc.UpdateUserInput{
		EditorID: 10,
		TargetID: 10,
		Name:     ptr("  Alice  "),
		Email:    ptr("  AlIce@Example.Com "),
	}

	match := mock.MatchedBy(func(f *repository.UserUpdateFields) bool {
		if f == nil || f.Name == nil || f.Email == nil {
			return false
		}
		return *f.Name == "Alice" &&
			*f.Email == "alice@example.com" &&
			f.PasswordHash == nil && f.SetAdmin == nil
	})
	updated := &domain.User{
		ID: 10, Name: "Alice", Email: ptr("alice@example.com"),
		IsAdmin: false, IsRoot: false, IsTest: false,
	}
	ur.EXPECT().UpdatePartial(ctx, 10, match).Return(updated, nil).Once()

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.NoError(t, err)
	require.Equal(t, 10, got.ID)
	require.Equal(t, "Alice", got.Name)
	require.NotNil(t, got.Email)
	require.Equal(t, "alice@example.com", *got.Email)
	require.True(t, commitCalled)
}

func TestUpdateUser_RootUpdatesRoleToAdmin_Success(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	editor := &domain.User{ID: 1, IsRoot: true}
	target := &domain.User{ID: 2, IsRoot: false, IsTest: false, IsAdmin: false}

	ur.EXPECT().FindForUpdate(ctx, 1).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 2).Return(target, nil).Once()

	in := uc.UpdateUserInput{EditorID: 1, TargetID: 2, Role: ptr("admin")}

	match := mock.MatchedBy(func(f *repository.UserUpdateFields) bool {
		return f != nil && f.SetAdmin != nil && *f.SetAdmin == true &&
			f.Name == nil && f.Email == nil && f.PasswordHash == nil
	})
	after := *target
	after.IsAdmin = true
	ur.EXPECT().UpdatePartial(ctx, 2, match).Return(&after, nil).Once()

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.NoError(t, err)
	require.True(t, got.IsAdmin)
	require.True(t, commitCalled)
}

func TestUpdateUser_SelfChangePassword_Success(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	oldHash, _ := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.DefaultCost)
	editor := &domain.User{ID: 10, IsRoot: false, IsTest: false}
	target := &domain.User{ID: 10, Password: string(oldHash)}

	ur.EXPECT().FindForUpdate(ctx, 10).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 10).Return(target, nil).Once()

	in := uc.UpdateUserInput{
		EditorID:        10,
		TargetID:        10,
		PasswordCurrent: ptr("oldpass"),
		PasswordNew:     ptr("newpass"),
	}

	match := mock.MatchedBy(func(f *repository.UserUpdateFields) bool {
		if f == nil || f.PasswordHash == nil {
			return false
		}
		return bcrypt.CompareHashAndPassword([]byte(*f.PasswordHash), []byte("newpass")) == nil &&
			f.Name == nil && f.Email == nil && f.SetAdmin == nil
	})

	after := &domain.User{ID: 10}
	ur.EXPECT().UpdatePartial(ctx, 10, match).Return(after, nil).Once()

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.NoError(t, err)
	require.Equal(t, 10, got.ID)
	require.True(t, commitCalled)
}

func TestUpdateUser_NonRootCannotUpdateOther(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	editor := &domain.User{ID: 10, IsRoot: false}
	target := &domain.User{ID: 20, IsRoot: false}

	ur.EXPECT().FindForUpdate(ctx, 10).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 20).Return(target, nil).Once()

	in := uc.UpdateUserInput{EditorID: 10, TargetID: 20, Name: ptr("X")}

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "Unauthorized")
	ur.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
	require.False(t, commitCalled)
}

func TestUpdateUser_TestUserCannotUpdateSelf(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	editor := &domain.User{ID: 10, IsRoot: false, IsTest: true}
	target := &domain.User{ID: 10, IsRoot: false, IsTest: true}

	ur.EXPECT().FindForUpdate(ctx, 10).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 10).Return(target, nil).Once()

	in := uc.UpdateUserInput{EditorID: 10, TargetID: 10, Name: ptr("X")}

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "Unauthorized")
	ur.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
	require.False(t, commitCalled)
}

func TestUpdateUser_InvalidEmail_ValidationError(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	editor := &domain.User{ID: 10}
	target := &domain.User{ID: 10}

	ur.EXPECT().FindForUpdate(ctx, 10).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 10).Return(target, nil).Once()

	in := uc.UpdateUserInput{EditorID: 10, TargetID: 10, Email: ptr("bad@@")}

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "VALIDATION")
	ur.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
	require.False(t, commitCalled)
}

func TestUpdateUser_InvalidRole_ValidationError(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	editor := &domain.User{ID: 1, IsRoot: true}
	target := &domain.User{ID: 2, IsRoot: false, IsTest: false}

	ur.EXPECT().FindForUpdate(ctx, 1).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 2).Return(target, nil).Once()

	in := uc.UpdateUserInput{EditorID: 1, TargetID: 2, Role: ptr("superadmin")}

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "VALIDATION")
	ur.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
	require.False(t, commitCalled)
}

func TestUpdateUser_ChangeRole_NonRootEditor_Unauthorized(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	editor := domain.User{ID: 3, IsRoot: false}
	target := domain.User{ID: 4, IsRoot: false, IsTest: false}
	ur.EXPECT().FindForUpdate(ctx, editor.ID).Return(&editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, target.ID).Return(&target, nil).Once()

	in := uc.UpdateUserInput{EditorID: editor.ID, TargetID: target.ID, Role: ptr("admin")}

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "Unauthorized")
	ur.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
	require.False(t, commitCalled)
}

func TestUpdateUser_ChangeRole_TargetIsRoot_Unauthorized(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	editor := domain.User{ID: 1, IsRoot: true}
	target := domain.User{ID: 2, IsRoot: true}
	ur.EXPECT().FindForUpdate(ctx, editor.ID).Return(&editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, target.ID).Return(&target, nil).Once()

	in := uc.UpdateUserInput{EditorID: editor.ID, TargetID: target.ID, Role: ptr("admin")}

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "Unauthorized")
	ur.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
	require.False(t, commitCalled)
}

func TestUpdateUser_ChangeRole_TargetIsTest_Unauthorized(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	editor := domain.User{ID: 1, IsRoot: true}
	target := domain.User{ID: 2, IsRoot: false, IsTest: true}
	ur.EXPECT().FindForUpdate(ctx, editor.ID).Return(&editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, target.ID).Return(&target, nil).Once()

	in := uc.UpdateUserInput{EditorID: editor.ID, TargetID: target.ID, Role: ptr("admin")}

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "Unauthorized")
	ur.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
	require.False(t, commitCalled)
}

func TestUpdateUser_PasswordChange_MissingCurrent_ValidationError(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	oldHash, _ := bcrypt.GenerateFromPassword([]byte("old"), bcrypt.DefaultCost)
	editor := &domain.User{ID: 5}
	target := &domain.User{ID: 5, Password: string(oldHash)}

	ur.EXPECT().FindForUpdate(ctx, 5).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 5).Return(target, nil).Once()

	in := uc.UpdateUserInput{EditorID: 5, TargetID: 5, PasswordNew: ptr("new")}

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "VALIDATION")
	ur.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
	require.False(t, commitCalled)
}

func TestUpdateUser_PasswordChange_WrongCurrent_InvalidCredential(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	oldHash, _ := bcrypt.GenerateFromPassword([]byte("old"), bcrypt.DefaultCost)
	editor := &domain.User{ID: 6}
	target := &domain.User{ID: 6, Password: string(oldHash)}

	ur.EXPECT().FindForUpdate(ctx, 6).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 6).Return(target, nil).Once()

	in := uc.UpdateUserInput{EditorID: 6, TargetID: 6, PasswordCurrent: ptr("wrong"), PasswordNew: ptr("new")}

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "ERR_INVALID_CREDENTIAL")
	ur.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
	require.False(t, commitCalled)
}

func TestUpdateUser_UpdatePartial_RepositoryError(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error { commitCalled = commit; return nil })

	editor := &domain.User{ID: 7}
	target := &domain.User{ID: 7}

	ur.EXPECT().FindForUpdate(ctx, 7).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 7).Return(target, nil).Once()

	in := uc.UpdateUserInput{EditorID: 7, TargetID: 7, Email: ptr("dup@example.com")}

	match := mock.MatchedBy(func(f *repository.UserUpdateFields) bool {
		return f != nil && f.Email != nil && *f.Email == "dup@example.com"
	})
	ur.EXPECT().UpdatePartial(ctx, 7, match).Return(nil, errors.New("conflict")).Once()

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "conflict")
	require.False(t, commitCalled)
}

func TestUpdateUser_TxBeginError(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	beginErr(ctx, tm, errors.New("tx-begin"))

	in := uc.UpdateUserInput{EditorID: 1, TargetID: 1}
	sut := newUCWithMocks(t, ur, tm)

	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "tx-begin")
	ur.AssertNotCalled(t, "FindForUpdate", mock.Anything, mock.Anything)
}

func TestUpdateUser_CommitError(t *testing.T) {
	ctx := context.Background()
	ur := usermock.NewMockRepository(t)
	tm := txmock.NewMockManager(t)
	commitCalled := false
	beginOK(ctx, tm, func(commit bool) error {
		commitCalled = commit
		if commit {
			return errors.New("commit-fail")
		}
		return nil
	})

	editor := &domain.User{ID: 8}
	target := &domain.User{ID: 8}

	ur.EXPECT().FindForUpdate(ctx, 8).Return(editor, nil).Once()
	ur.EXPECT().FindForUpdate(ctx, 8).Return(target, nil).Once()

	in := uc.UpdateUserInput{EditorID: 8, TargetID: 8, Name: ptr("Z")}

	match := mock.MatchedBy(func(f *repository.UserUpdateFields) bool {
		return f != nil && f.Name != nil && *f.Name == "Z"
	})
	after := &domain.User{ID: 8, Name: "Z"}
	ur.EXPECT().UpdatePartial(ctx, 8, match).Return(after, nil).Once()

	sut := newUCWithMocks(t, ur, tm)
	got, err := sut.UpdateUser(ctx, in)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "commit-fail")
	require.True(t, commitCalled)
}
