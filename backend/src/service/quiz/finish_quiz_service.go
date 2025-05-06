package quiz_service

import (
	"context"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/quizquestion"
	"word_app/backend/src/models"
)

func (s *QuizServiceImpl) finishQuizTx(
	ctx context.Context,
	tx *ent.Tx,
	q *ent.Quiz,
) (*models.Result, error) {

	// 集計はメモリ側で
	qs := q.Edges.QuizQuestions
	var (
		correct         int
		resultQuestions = make([]models.ResultQuestion, 0, len(q.Edges.QuizQuestions))
	)

	for _, x := range qs {
		// nil セーフガード
		isCor := x.IsCorrect != nil && *x.IsCorrect
		if isCor {
			correct++
		}

		resQ := models.ResultQuestion{
			QuestionNumber: x.QuestionNumber,
			WordName:       x.WordName,
			PosID:          x.PosID,
			CorrectJpmId:   x.CorrectJpmID,
			ChoicesJpms:    x.ChoicesJpms,
			AnswerJpmId:    derefInt(x.AnswerJpmID),
			IsCorrect:      isCor,
			TimeMs:         derefInt(x.TimeMs),
		}
		resultQuestions = append(resultQuestions, resQ)
	}
	rate := float64(correct) / float64(q.TotalQuestionsCount)

	if _, err := tx.Quiz.
		UpdateOne(q).
		SetIsRunning(false).
		SetCorrectCount(correct).
		SetResultCorrectRate(rate).
		Save(ctx); err != nil {
		return nil, err
	}

	if !q.IsSaveResult {
		now := time.Now()
		if _, err := tx.Quiz.UpdateOne(q).SetDeletedAt(now).Save(ctx); err != nil {
			return nil, err
		}
		if _, err := tx.QuizQuestion.
			Update().
			Where(quizquestion.QuizIDEQ(q.ID)).
			SetDeletedAt(now).
			Save(ctx); err != nil {
			return nil, err
		}
	}

	// ─── 5. レスポンス組立て ───────────────────────────────
	return &models.Result{
		QuizNumber:          q.QuizNumber,
		TotalQuestionsCount: q.TotalQuestionsCount,
		CorrectCount:        correct,
		ResultCorrectRate:   rate,
		ResultSetting: models.ResultSetting{
			IsSaveResult:        q.IsSaveResult,
			IsRegisteredWords:   q.IsRegisteredWords,
			SettingCorrectRate:  q.SettingCorrectRate,
			IsIdioms:            q.IsIdioms,
			IsSpecialCharacters: q.IsSpecialCharacters,
			AttentionLevelList:  q.AttentionLevelList,
			ChoicesPosIds:       q.ChoicesPosIds,
		},
		ResultQuestions: resultQuestions,
	}, nil
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
