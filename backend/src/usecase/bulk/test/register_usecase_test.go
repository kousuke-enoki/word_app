package bulk_test

import (
	"context"
	"errors"
	"testing"

	"word_app/backend/config"
	regmock "word_app/backend/src/mocks/infrastructure/repository/registeredword"
	txmock "word_app/backend/src/mocks/infrastructure/repository/tx"
	usermock "word_app/backend/src/mocks/infrastructure/repository/user"
	wordmock "word_app/backend/src/mocks/infrastructure/repository/word"
	"word_app/backend/src/usecase/bulk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeRegisterUC(t *testing.T, wordRepo *wordmock.MockReadRepository, rwRead *regmock.MockReadRepository, rwWrite *regmock.MockWriteRepository, tm *txmock.MockManager, userRepo *usermock.MockRepository, limits *config.LimitsCfg) bulk.RegisterUsecase {
	return bulk.NewRegisterUsecase(wordRepo, rwRead, rwWrite, tm, userRepo, limits)
}

func TestRegisterUsecase_Register(t *testing.T) {
	ctx := context.Background()
	limits := &config.LimitsCfg{
		RegisteredWordsPerUser: 200,
		BulkRegisterMaxItems:   200,
	}

	t.Run("success - all words registered", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		// FindIDsByNamesはBeginより前なのでctxで呼ばれる
		wordRepo.On("FindIDsByNames", ctx, []string{"apple", "banana"}).
			Return(map[string]int{"apple": 1, "banana": 2}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("LockByID", ctx, 1).Return(nil)
		rwRead.On("CountActiveByUser", ctx, 1).Return(0, nil)
		// mapの順序は不定なので、mock.MatchedByを使用
		rwRead.On("FindActiveMapByUserAndWordIDs", ctx, 1, mock.MatchedBy(func(ids []int) bool {
			idMap := make(map[int]bool)
			for _, id := range ids {
				idMap[id] = true
			}
			return len(ids) == 2 && idMap[1] && idMap[2]
		})).Return(map[int]bool{}, nil)
		// 順序不定なのでmock.Anythingを使用
		rwWrite.On("CreateActive", ctx, 1, mock.Anything).Return(nil).Twice()

		result, err := uc.Register(ctx, 1, []string{"apple", "banana"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []string{"apple", "banana"}, result.Success)
		assert.Empty(t, result.Failed)
		tm.AssertExpectations(t)
	})

	t.Run("success - partial success with some failed", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		wordRepo.On("FindIDsByNames", ctx, []string{"apple", "notexists"}).
			Return(map[string]int{"apple": 1}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("LockByID", ctx, 1).Return(nil)
		rwRead.On("CountActiveByUser", ctx, 1).Return(0, nil)
		rwRead.On("FindActiveMapByUserAndWordIDs", ctx, 1, []int{1}).
			Return(map[int]bool{}, nil)
		rwWrite.On("CreateActive", ctx, 1, 1).Return(nil)

		result, err := uc.Register(ctx, 1, []string{"apple", "notexists"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []string{"apple"}, result.Success)
		assert.Len(t, result.Failed, 1)
		assert.Equal(t, "notexists", result.Failed[0].Word)
		assert.Equal(t, "not_exists", result.Failed[0].Reason)
		tm.AssertExpectations(t)
	})

	t.Run("success - already registered", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		wordRepo.On("FindIDsByNames", ctx, []string{"apple"}).
			Return(map[string]int{"apple": 1}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("LockByID", ctx, 1).Return(nil)
		rwRead.On("CountActiveByUser", ctx, 1).Return(0, nil)
		rwRead.On("FindActiveMapByUserAndWordIDs", ctx, 1, []int{1}).
			Return(map[int]bool{1: true}, nil)

		result, err := uc.Register(ctx, 1, []string{"apple"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Success)
		assert.Len(t, result.Failed, 1)
		assert.Equal(t, "apple", result.Failed[0].Word)
		assert.Equal(t, "already_registered", result.Failed[0].Reason)
		tm.AssertExpectations(t)
	})

	t.Run("success - reactivate inactive word", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		wordRepo.On("FindIDsByNames", ctx, []string{"apple"}).
			Return(map[string]int{"apple": 1}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("LockByID", ctx, 1).Return(nil)
		rwRead.On("CountActiveByUser", ctx, 1).Return(0, nil)
		rwRead.On("FindActiveMapByUserAndWordIDs", ctx, 1, []int{1}).
			Return(map[int]bool{1: false}, nil)
		rwWrite.On("Activate", ctx, 1, 1).Return(nil)

		result, err := uc.Register(ctx, 1, []string{"apple"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []string{"apple"}, result.Success)
		assert.Empty(t, result.Failed)
		tm.AssertExpectations(t)
	})

	t.Run("error - empty payload", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		result, err := uc.Register(ctx, 1, []string{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty payload")
		assert.Nil(t, result)
		tm.AssertNotCalled(t, "Begin", mock.Anything)
	})

	t.Run("error - too many words", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		words := make([]string, 201)
		for i := range words {
			words[i] = "word"
		}

		result, err := uc.Register(ctx, 1, words)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too many words")
		assert.Nil(t, result)
		tm.AssertNotCalled(t, "Begin", mock.Anything)
	})

	t.Run("error - limit reached (all failed)", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		wordRepo.On("FindIDsByNames", ctx, []string{"apple", "banana"}).
			Return(map[string]int{"apple": 1, "banana": 2}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("LockByID", ctx, 1).Return(nil)
		rwRead.On("CountActiveByUser", ctx, 1).Return(200, nil)

		result, err := uc.Register(ctx, 1, []string{"apple", "banana"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Nil(t, result.Success)
		assert.Len(t, result.Failed, 2)
		assert.Equal(t, "apple", result.Failed[0].Word)
		assert.Equal(t, "limit_reached", result.Failed[0].Reason)
		assert.Equal(t, "banana", result.Failed[1].Word)
		assert.Equal(t, "limit_reached", result.Failed[1].Reason)
		tm.AssertExpectations(t)
	})

	t.Run("error - FindIDsByNames fails", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		wordRepo.On("FindIDsByNames", ctx, []string{"apple"}).
			Return(nil, errors.New("database error"))

		result, err := uc.Register(ctx, 1, []string{"apple"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.Nil(t, result)
		tm.AssertExpectations(t)
	})

	t.Run("error - CountActiveByUser fails", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		wordRepo.On("FindIDsByNames", ctx, []string{"apple"}).
			Return(map[string]int{"apple": 1}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("LockByID", ctx, 1).Return(nil)
		rwRead.On("CountActiveByUser", ctx, 1).Return(0, errors.New("database error"))

		result, err := uc.Register(ctx, 1, []string{"apple"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.Nil(t, result)
		tm.AssertExpectations(t)
	})

	t.Run("error - CreateActive fails", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		wordRepo.On("FindIDsByNames", ctx, []string{"apple"}).
			Return(map[string]int{"apple": 1}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("LockByID", ctx, 1).Return(nil)
		rwRead.On("CountActiveByUser", ctx, 1).Return(0, nil)
		rwRead.On("FindActiveMapByUserAndWordIDs", ctx, 1, []int{1}).
			Return(map[int]bool{}, nil)
		rwWrite.On("CreateActive", ctx, 1, 1).Return(errors.New("database error"))

		result, err := uc.Register(ctx, 1, []string{"apple"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Success)
		assert.Len(t, result.Failed, 1)
		assert.Equal(t, "apple", result.Failed[0].Word)
		assert.Equal(t, "db_error", result.Failed[0].Reason)
		tm.AssertExpectations(t)
	})

	t.Run("success - normalize case and duplicates", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		wordRepo.On("FindIDsByNames", ctx, []string{"apple"}).
			Return(map[string]int{"apple": 1}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("LockByID", ctx, 1).Return(nil)
		rwRead.On("CountActiveByUser", ctx, 1).Return(0, nil)
		rwRead.On("FindActiveMapByUserAndWordIDs", ctx, 1, []int{1}).
			Return(map[int]bool{}, nil)
		rwWrite.On("CreateActive", ctx, 1, 1).Return(nil)

		result, err := uc.Register(ctx, 1, []string{"Apple", "APPLE", "apple"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []string{"apple"}, result.Success)
		assert.Empty(t, result.Failed)
		tm.AssertExpectations(t)
	})

	t.Run("success - mixed scenarios", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		rwRead := regmock.NewMockReadRepository(t)
		rwWrite := regmock.NewMockWriteRepository(t)
		tm := txmock.NewMockManager(t)
		userRepo := usermock.NewMockRepository(t)

		uc := makeRegisterUC(t, wordRepo, rwRead, rwWrite, tm, userRepo, limits)

		wordRepo.On("FindIDsByNames", ctx, []string{"apple", "banana", "notexists", "already"}).
			Return(map[string]int{"apple": 1, "banana": 2, "already": 3}, nil)
		tm.On("Begin", ctx).Return(ctx, func(bool) error { return nil }, nil)
		userRepo.On("LockByID", ctx, 1).Return(nil)
		rwRead.On("CountActiveByUser", ctx, 1).Return(100, nil)
		// mapの順序は不定なので、mock.MatchedByを使用
		rwRead.On("FindActiveMapByUserAndWordIDs", ctx, 1, mock.MatchedBy(func(ids []int) bool {
			idMap := make(map[int]bool)
			for _, id := range ids {
				idMap[id] = true
			}
			return len(ids) == 3 && idMap[1] && idMap[2] && idMap[3]
		})).Return(map[int]bool{3: true}, nil)
		rwWrite.On("CreateActive", ctx, 1, 1).Return(nil)
		rwWrite.On("CreateActive", ctx, 1, 2).Return(nil)

		result, err := uc.Register(ctx, 1, []string{"apple", "banana", "notexists", "already"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []string{"apple", "banana"}, result.Success)
		assert.Len(t, result.Failed, 2)
		assert.Equal(t, "notexists", result.Failed[0].Word)
		assert.Equal(t, "not_exists", result.Failed[0].Reason)
		assert.Equal(t, "already", result.Failed[1].Word)
		assert.Equal(t, "already_registered", result.Failed[1].Reason)
		tm.AssertExpectations(t)
	})
}
