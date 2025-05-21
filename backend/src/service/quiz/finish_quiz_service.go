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

	"github.com/sirupsen/logrus"
)

func (s *QuizServiceImpl) finishQuizTx(
	ctx context.Context,
	tx *ent.Tx,
	q *ent.Quiz,
	userID int,
) error {
	logrus.Info("finishQuizTx")
	logrus.Info(q)
	logrus.Info(q.Edges)
	// 集計はメモリ側で

	// -------- ensure questions are loaded --------
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
		// resultQuestions = make([]models.ResultQuestion, 0, len(q.Edges.QuizQuestions))
	)

	for _, x := range qs {
		logrus.Info("x")
		logrus.Info(x)
		// var (
		// 	IsRegistered   bool
		// 	AttentionLevel int
		// 	QuizCount      int
		// 	CorrectCount   int
		// 	CorrectRate    int
		// )
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
			// if ent.IsNotFound(err) {
			// 	return nil, err
			// }
			// if err != nil {
			// 	return nil, errors.New("failed to query RegisteredWord")
			// }
		switch {
		case ent.IsNotFound(err): // --- 存在しない → 新規作成 ---
			// IsRegistered = false
			// AttentionLevel = 1
			// QuizCount = 1
			// if isCor {
			// 	CorrectCount++
			// }
			// CorrectRate = CorrectCount * 100 / QuizCount
			initCor := 0
			if isCor {
				initCor = 1
			}
			_, err = tx.RegisteredWord.
				Create().
				SetUserID(userID).
				SetWordID(x.WordID).
				SetIsActive(false).   // 初期値
				SetAttentionLevel(1). // 初期値
				SetQuizCount(1).
				SetCorrectCount(initCor).
				SetCorrectRate(initCor * 100).
				Save(ctx)
			if err != nil {
				return errors.New("Failed to create RegisteredWord")
			}

		case err != nil: // --- DB エラー ---
			return fmt.Errorf("query registered_word: %w", err)

		default: // 既存レコードをインクリメント
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

	// 	resQ := models.ResultQuestion{
	// 		QuestionNumber: x.QuestionNumber,
	// 		WordName:       x.WordName,
	// 		PosID:          x.PosID,
	// 		CorrectJpmId:   x.CorrectJpmID,
	// 		ChoicesJpms:    x.ChoicesJpms,
	// 		AnswerJpmId:    derefInt(x.AnswerJpmID),
	// 		IsCorrect:      isCor,
	// 		TimeMs:         derefInt(x.TimeMs),
	// 		RegisteredWord: models.RegisteredWord{
	// 			IsRegistered:   IsRegistered,
	// 			AttentionLevel: AttentionLevel,
	// 			QuizCount:      QuizCount,
	// 			CorrectCount:   CorrectCount,
	// 		},
	// 	}
	// 	logrus.Info("resQ")
	// 	logrus.Info(resQ)
	// 	resultQuestions = append(resultQuestions, resQ)
	// }
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
	// if !q.IsSaveResult {
	// 	now := time.Now()
	// 	if _, err := tx.Quiz.UpdateOneID(q.ID).SetDeletedAt(now).Save(ctx); err != nil {
	// 		return err
	// 	}
	// 	if _, err := tx.QuizQuestion.
	// 		Update().
	// 		Where(quizquestion.QuizIDEQ(q.ID)).
	// 		SetDeletedAt(now).
	// 		Save(ctx); err != nil {
	// 		return err
	// 	}
	// }
	// logrus.Info("resultQuestions")
	// logrus.Info(resultQuestions)
	// // ─── 5. レスポンス組立て ───────────────────────────────
	// return &models.Result{
	// 	QuizNumber:          q.QuizNumber,
	// 	TotalQuestionsCount: q.TotalQuestionsCount,
	// 	CorrectCount:        correct,
	// 	ResultCorrectRate:   rate,
	// 	ResultSetting: models.ResultSetting{
	// 		IsSaveResult:        q.IsSaveResult,
	// 		IsRegisteredWords:   q.IsRegisteredWords,
	// 		SettingCorrectRate:  q.SettingCorrectRate,
	// 		IsIdioms:            q.IsIdioms,
	// 		IsSpecialCharacters: q.IsSpecialCharacters,
	// 		AttentionLevelList:  q.AttentionLevelList,
	// 		ChoicesPosIds:       q.ChoicesPosIds,
	// 	},
	// 	ResultQuestions: resultQuestions,
	// }, nil
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
