package bulk_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"word_app/backend/config"
	"word_app/backend/src/domain"
	regmock "word_app/backend/src/mocks/infrastructure/repository/registeredword"
	udumock "word_app/backend/src/mocks/infrastructure/repository/userdailyusage"
	wordmock "word_app/backend/src/mocks/infrastructure/repository/word"
	"word_app/backend/src/usecase/bulk"
	"word_app/backend/src/usecase/shared/ucerr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClock struct {
	now time.Time
}

func (m *mockClock) Now() time.Time {
	return m.now
}

func makeTokenizeUC(t *testing.T, wordRepo *wordmock.MockReadRepository, regReadRepo *regmock.MockReadRepository, userDailyUsageRepo *udumock.MockRepository, clock *mockClock, limits *config.LimitsCfg) bulk.TokenizeUsecase {
	return bulk.NewTokenizeUsecase(wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)
}

func TestTokenizeUsecase_Execute(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	clock := &mockClock{now: now}
	limits := &config.LimitsCfg{
		BulkMaxPerDay:         5,
		BulkTokenizeMaxTokens: 200,
	}

	// 各テストケースを独立したサブテストとして実行
	// これによりサイクロマティック複雑度を下げる
	t.Run("success cases", func(t *testing.T) {
		testTokenizeExecute_SuccessCases(t, ctx, now, clock, limits)
	})

	t.Run("error cases", func(t *testing.T) {
		testTokenizeExecute_ErrorCases(t, ctx, now, clock, limits)
	})
}

func testTokenizeExecute_SuccessCases(t *testing.T, ctx context.Context, now time.Time, clock *mockClock, limits *config.LimitsCfg) {
	t.Run("candidates and registered words", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		userDailyUsageRepo.On("IncBulkOr429", ctx, 1, now, 999999).
			Return(&domain.DailyUsageUpdateResult{BulkCount: 1}, nil)
		wordRepo.On("FindIDsByNames", ctx, []string{"hello", "world", "test"}).
			Return(map[string]int{"hello": 1, "world": 2, "test": 3}, nil)
		// mapの順序は不定なので、mock.MatchedByを使用
		regReadRepo.On("ActiveWordIDSetByUser", ctx, 1, mock.MatchedBy(func(ids []int) bool {
			idMap := make(map[int]bool)
			for _, id := range ids {
				idMap[id] = true
			}
			return len(ids) == 3 && idMap[1] && idMap[2] && idMap[3]
		})).Return(map[int]struct{}{1: {}}, nil)

		cands, regs, notExist, err := uc.Execute(ctx, 1, "Hello world test")

		assert.NoError(t, err)
		assert.Equal(t, []string{"world", "test"}, cands)
		assert.Equal(t, []string{"hello"}, regs)
		assert.Empty(t, notExist)
		userDailyUsageRepo.AssertExpectations(t)
	})

	t.Run("all candidates", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		userDailyUsageRepo.On("IncBulkOr429", ctx, 1, now, 999999).
			Return(&domain.DailyUsageUpdateResult{BulkCount: 1}, nil)
		wordRepo.On("FindIDsByNames", ctx, []string{"hello", "world"}).
			Return(map[string]int{"hello": 1, "world": 2}, nil)
		// mapの順序は不定なので、mock.MatchedByを使用
		regReadRepo.On("ActiveWordIDSetByUser", ctx, 1, mock.MatchedBy(func(ids []int) bool {
			idMap := make(map[int]bool)
			for _, id := range ids {
				idMap[id] = true
			}
			return len(ids) == 2 && idMap[1] && idMap[2]
		})).Return(map[int]struct{}{}, nil)

		cands, regs, notExist, err := uc.Execute(ctx, 1, "Hello world")

		assert.NoError(t, err)
		assert.Equal(t, []string{"hello", "world"}, cands)
		assert.Empty(t, regs)
		assert.Empty(t, notExist)
		userDailyUsageRepo.AssertExpectations(t)
	})

	t.Run("with not exist words", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		userDailyUsageRepo.On("IncBulkOr429", ctx, 1, now, 999999).
			Return(&domain.DailyUsageUpdateResult{BulkCount: 1}, nil)
		wordRepo.On("FindIDsByNames", ctx, []string{"hello", "notexists"}).
			Return(map[string]int{"hello": 1}, nil)
		regReadRepo.On("ActiveWordIDSetByUser", ctx, 1, mock.MatchedBy(func(ids []int) bool {
			return len(ids) == 1 && ids[0] == 1
		})).Return(map[int]struct{}{}, nil)

		cands, regs, notExist, err := uc.Execute(ctx, 1, "Hello notexists")

		assert.NoError(t, err)
		assert.Equal(t, []string{"hello"}, cands)
		assert.Empty(t, regs)
		assert.Equal(t, []string{"notexists"}, notExist)
		userDailyUsageRepo.AssertExpectations(t)
	})

	t.Run("empty text", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		userDailyUsageRepo.On("IncBulkOr429", ctx, 1, now, 999999).
			Return(&domain.DailyUsageUpdateResult{BulkCount: 1}, nil)

		cands, regs, notExist, err := uc.Execute(ctx, 1, "")

		assert.NoError(t, err)
		assert.Nil(t, cands)
		assert.Nil(t, regs)
		assert.Nil(t, notExist)
		userDailyUsageRepo.AssertExpectations(t)
	})

	t.Run("normalize case and duplicates", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		userDailyUsageRepo.On("IncBulkOr429", ctx, 1, now, 999999).
			Return(&domain.DailyUsageUpdateResult{BulkCount: 1}, nil)
		wordRepo.On("FindIDsByNames", ctx, []string{"hello"}).
			Return(map[string]int{"hello": 1}, nil)
		regReadRepo.On("ActiveWordIDSetByUser", ctx, 1, mock.MatchedBy(func(ids []int) bool {
			return len(ids) == 1 && ids[0] == 1
		})).Return(map[int]struct{}{}, nil)

		cands, regs, notExist, err := uc.Execute(ctx, 1, "Hello HELLO hello")

		assert.NoError(t, err)
		assert.Equal(t, []string{"hello"}, cands) // 正規化されて1つ
		assert.Empty(t, regs)
		assert.Empty(t, notExist)
		userDailyUsageRepo.AssertExpectations(t)
	})

	t.Run("apostrophes", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		userDailyUsageRepo.On("IncBulkOr429", ctx, 1, now, 999999).
			Return(&domain.DailyUsageUpdateResult{BulkCount: 1}, nil)
		wordRepo.On("FindIDsByNames", ctx, []string{"don't", "can't"}).
			Return(map[string]int{"don't": 1, "can't": 2}, nil)
		// mapの順序は不定なので、mock.MatchedByを使用
		regReadRepo.On("ActiveWordIDSetByUser", ctx, 1, mock.MatchedBy(func(ids []int) bool {
			idMap := make(map[int]bool)
			for _, id := range ids {
				idMap[id] = true
			}
			return len(ids) == 2 && idMap[1] && idMap[2]
		})).Return(map[int]struct{}{}, nil)

		cands, regs, notExist, err := uc.Execute(ctx, 1, "don't can't")

		assert.NoError(t, err)
		assert.Equal(t, []string{"don't", "can't"}, cands)
		assert.Empty(t, regs)
		assert.Empty(t, notExist)
		userDailyUsageRepo.AssertExpectations(t)
	})
	t.Run("non-test user unlimited quota", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		userDailyUsageRepo.On("IncBulkOr429", ctx, 2, now, 999999).
			Return(&domain.DailyUsageUpdateResult{BulkCount: 999999}, nil)
		wordRepo.On("FindIDsByNames", ctx, []string{"hello"}).
			Return(map[string]int{"hello": 1}, nil)
		regReadRepo.On("ActiveWordIDSetByUser", ctx, 2, mock.MatchedBy(func(ids []int) bool {
			return len(ids) == 1 && ids[0] == 1
		})).Return(map[int]struct{}{}, nil)

		cands, regs, notExist, err := uc.Execute(ctx, 2, "Hello")

		assert.NoError(t, err)
		assert.Equal(t, []string{"hello"}, cands)
		assert.Empty(t, regs)
		assert.Empty(t, notExist)
		// 999999の制限で呼ばれていることを確認
		userDailyUsageRepo.AssertExpectations(t)
	})
}

func testTokenizeExecute_ErrorCases(t *testing.T, ctx context.Context, now time.Time, clock *mockClock, limits *config.LimitsCfg) {
	t.Run("quota exceeded", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		userDailyUsageRepo.On("IncBulkOr429", ctx, 1, now, 999999).
			Return(nil, ucerr.TooManyRequests("daily quota exceeded"))

		cands, regs, notExist, err := uc.Execute(ctx, 1, "Hello world")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "quota")
		assert.Nil(t, cands)
		assert.Nil(t, regs)
		assert.Nil(t, notExist)
		userDailyUsageRepo.AssertExpectations(t)
	})

	t.Run("too many tokens", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		// 1001個のトークンを生成（上限200*5=1000を超える）
		largeText := ""
		for i := 0; i < 1001; i++ {
			largeText += "word "
		}

		userDailyUsageRepo.On("IncBulkOr429", ctx, 1, now, 999999).
			Return(&domain.DailyUsageUpdateResult{BulkCount: 1}, nil)

		cands, regs, notExist, err := uc.Execute(ctx, 1, largeText)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too many tokens")
		assert.Nil(t, cands)
		assert.Nil(t, regs)
		assert.Nil(t, notExist)
		userDailyUsageRepo.AssertExpectations(t)
	})

	t.Run("FindIDsByNames fails", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		userDailyUsageRepo.On("IncBulkOr429", ctx, 1, now, 999999).
			Return(&domain.DailyUsageUpdateResult{BulkCount: 1}, nil)
		wordRepo.On("FindIDsByNames", ctx, []string{"hello"}).
			Return(nil, errors.New("database error"))

		cands, regs, notExist, err := uc.Execute(ctx, 1, "Hello")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.Nil(t, cands)
		assert.Nil(t, regs)
		assert.Nil(t, notExist)
		userDailyUsageRepo.AssertExpectations(t)
	})

	t.Run("ActiveWordIDSetByUser fails", func(t *testing.T) {
		wordRepo := wordmock.NewMockReadRepository(t)
		regReadRepo := regmock.NewMockReadRepository(t)
		userDailyUsageRepo := udumock.NewMockRepository(t)

		uc := makeTokenizeUC(t, wordRepo, regReadRepo, userDailyUsageRepo, clock, limits)

		userDailyUsageRepo.On("IncBulkOr429", ctx, 1, now, 999999).
			Return(&domain.DailyUsageUpdateResult{BulkCount: 1}, nil)
		wordRepo.On("FindIDsByNames", ctx, []string{"hello"}).
			Return(map[string]int{"hello": 1}, nil)
		regReadRepo.On("ActiveWordIDSetByUser", ctx, 1, mock.MatchedBy(func(ids []int) bool {
			return len(ids) == 1 && ids[0] == 1
		})).Return(nil, errors.New("database error"))

		cands, regs, notExist, err := uc.Execute(ctx, 1, "Hello")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.Nil(t, cands)
		assert.Nil(t, regs)
		assert.Nil(t, notExist)
		userDailyUsageRepo.AssertExpectations(t)
	})
}
