package quiz_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/quiz"
	"word_app/backend/ent/quizquestion"
	"word_app/backend/ent/registeredword"
)

// quiz終了処理
func (s *QuizServiceImpl) finishQuizTx(
	ctx context.Context,
	tx *ent.Tx,
	q *ent.Quiz,
	userID int,
) error {
	if q.Edges.QuizQuestions == nil {
		var err error
		q, err = tx.Quiz.
			Query().
			Where(quiz.IDEQ(q.ID)).
			WithQuizQuestions().
			Only(ctx)
		if err != nil {
			return err
		}
	}

	qs := q.Edges.QuizQuestions
	var (
		correct int
	)

	for _, x := range qs {
		// nil セーフガード
		isCor := x.IsCorrect != nil && *x.IsCorrect
		if isCor {
			correct++
		}

		// registeredWord 取得、更新
		rw, err := tx.RegisteredWord.
			Query().
			Where(
				registeredword.WordID(x.WordID),
				registeredword.UserID(userID),
			).
			Only(ctx)
		switch {
		case ent.IsNotFound(err): // --- 存在しない → 新規作成 ---
			initCor := 0
			if isCor {
				initCor = 1
			}
			_, err = tx.RegisteredWord.
				Create().
				SetUserID(userID).
				SetWordID(x.WordID).
				SetIsActive(false).
				SetAttentionLevel(1).
				SetQuizCount(1).
				SetCorrectCount(initCor).
				SetCorrectRate(initCor * 100).
				Save(ctx)
			if err != nil {
				return errors.New("Failed to create RegisteredWord")
			}

		case err != nil:
			return fmt.Errorf("query registered_word: %w", err)

		default:
			newQuiz := rw.QuizCount + 1
			newCorrect := rw.CorrectCount
			if isCor {
				newCorrect++
			}
			_, err = tx.RegisteredWord.
				UpdateOneID(rw.ID).
				SetQuizCount(newQuiz).
				SetCorrectCount(newCorrect).
				SetCorrectRate(newCorrect * 100 / newQuiz).
				Save(ctx)
		}
		if err != nil {
			return errors.New("Failed to update RegisteredWord")
		}
	}
	// --- quiz 本体を完了状態に更新 ----------------------------
	rate := float64(correct) * 100 / float64(q.TotalQuestionsCount)

	if _, err := tx.Quiz.
		UpdateOneID(q.ID).
		SetIsRunning(false).
		SetCorrectCount(correct).
		SetResultCorrectRate(rate).
		Save(ctx); err != nil {
		return err
	}

	// --- 保存しない設定なら soft-delete -----------------------
	if !q.IsSaveResult {
		now := time.Now()
		if _, err := tx.Quiz.
			UpdateOneID(q.ID).
			SetDeletedAt(now).
			Save(ctx); err != nil {
			return err
		}
		if _, err := tx.QuizQuestion.
			Update().
			Where(quizquestion.QuizIDEQ(q.ID)).
			SetDeletedAt(now).
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func finishTx(perr *error, tx *ent.Tx) {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if *perr != nil {
		_ = tx.Rollback()
	} else {
		*perr = tx.Commit()
	}
}

// ヘルパ
func derefInt(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}
