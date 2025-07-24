package quiz

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

/*==================== public ====================*/

func (s *ServiceImpl) finishQuizTx(
	ctx context.Context,
	tx *ent.Tx,
	q *ent.Quiz,
	userID int,
) error {

	// ① 質問一覧を確実に取得
	q, err := s.loadQuizQuestions(ctx, tx, q)
	if err != nil {
		return err
	}

	// ② 各質問を処理し、正答数をカウント
	correct := 0
	for _, qq := range q.Edges.QuizQuestions {
		isCor, err := s.upsertRegisteredWord(ctx, tx, qq, userID)
		if err != nil {
			return err
		}
		if isCor {
			correct++
		}
	}

	// ③ クイズ本体を完了状態に更新
	if err := s.updateQuizResult(ctx, tx, q, correct); err != nil {
		return err
	}

	// ④ 保存しない設定なら soft-delete
	if !q.IsSaveResult {
		if err := s.softDeleteQuiz(ctx, tx, q.ID); err != nil {
			return err
		}
	}
	return nil
}

/*==================== repository-like helpers ====================*/

// Edge が無い場合のみ再取得
func (s *ServiceImpl) loadQuizQuestions(
	ctx context.Context,
	tx *ent.Tx,
	q *ent.Quiz,
) (*ent.Quiz, error) {
	if q.Edges.QuizQuestions != nil {
		return q, nil
	}
	return tx.Quiz.
		Query().
		Where(quiz.IDEQ(q.ID)).
		WithQuizQuestions().
		Only(ctx)
}

// RegisteredWord を作成 / 更新し、正答かどうかを返す
func (s *ServiceImpl) upsertRegisteredWord(
	ctx context.Context,
	tx *ent.Tx,
	qq *ent.QuizQuestion,
	userID int,
) (isCorrect bool, err error) {

	isCorrect = qq.IsCorrect != nil && *qq.IsCorrect

	rw, err := tx.RegisteredWord.
		Query().
		Where(
			registeredword.WordID(qq.WordID),
			registeredword.UserID(userID),
		).
		Only(ctx)

	switch {
	case ent.IsNotFound(err): // 新規作成
		initCor := 0
		if isCorrect {
			initCor = 1
		}
		_, err = tx.RegisteredWord.
			Create().
			SetUserID(userID).
			SetWordID(qq.WordID).
			SetIsActive(false).
			SetAttentionLevel(1).
			SetQuizCount(1).
			SetCorrectCount(initCor).
			SetCorrectRate(initCor * 100).
			Save(ctx)

	case err != nil: // クエリエラー
		return false, fmt.Errorf("query registered_word: %w", err)

	default: // 既存行を更新
		newQuiz := rw.QuizCount + 1
		newCorrect := rw.CorrectCount
		if isCorrect {
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
		return false, errors.New("failed to upsert RegisteredWord")
	}
	return
}

// クイズ本体の完了情報を更新
func (s *ServiceImpl) updateQuizResult(
	ctx context.Context,
	tx *ent.Tx,
	q *ent.Quiz,
	correct int,
) error {
	rate := float64(correct) * 100 / float64(q.TotalQuestionsCount)
	_, err := tx.Quiz.
		UpdateOneID(q.ID).
		SetIsRunning(false).
		SetCorrectCount(correct).
		SetResultCorrectRate(rate).
		Save(ctx)
	return err
}

// Quiz / QuizQuestion を soft-delete
func (s *ServiceImpl) softDeleteQuiz(
	ctx context.Context,
	tx *ent.Tx,
	quizID int,
) error {
	now := time.Now()
	if _, err := tx.Quiz.
		UpdateOneID(quizID).
		SetDeletedAt(now).
		Save(ctx); err != nil {
		return err
	}
	_, err := tx.QuizQuestion.
		Update().
		Where(quizquestion.QuizIDEQ(quizID)).
		SetDeletedAt(now).
		Save(ctx)
	return err
}

/*==================== 既存 util (変更なし) ====================*/

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
