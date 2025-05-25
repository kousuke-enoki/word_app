package quiz

import "context"

type WordRepo interface {
	RandomSelectableWords(ctx context.Context, userID int, filter WordFilter, limit int) ([]*Word, error)
	BeginTx(ctx context.Context) (Tx, error)
}

type Tx interface {
	CreateQuiz(ctx context.Context, userID int, in QuizRecord) (int /*quizID*/, error)
	CreateQuestions(ctx context.Context, quizID int, qs []QuizQuestionRecord) error
	Commit() error
	Rollback() error
}

// ───────── 単語取得用フィルタ ─────────
type WordFilter struct {
	PartsOfSpeech   []int
	RegisteredMode  int // 0:全 1:登録 2:未登録
	MaxCorrectRate  int
	AttentionLevels []int
	IncludeIdioms   int
	IncludeSpecial  int
}
