package quiz_service

import (
	"context"
	"errors"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/quizquestion"
	"word_app/backend/ent/registeredword"
	"word_app/backend/src/models"
)

func (s *QuizServiceImpl) finishQuizTx(
	ctx context.Context,
	tx *ent.Tx,
	q *ent.Quiz,
	userID int,
) (*models.Result, error) {

	// 集計はメモリ側で
	qs := q.Edges.QuizQuestions
	var (
		correct         int
		resultQuestions = make([]models.ResultQuestion, 0, len(q.Edges.QuizQuestions))
	)

	for _, x := range qs {
		var (
			IsRegistered   bool
			AttentionLevel int
			QuizCount      int
			CorrectCount   int
			CorrectRate    int
		)
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
		if ent.IsNotFound(err) {
			return nil, nil // 登録が存在しない場合はnilを返す
		}
		if err != nil {
			return nil, errors.New("failed to query RegisteredWord")
		}
		if rw == nil {
			IsRegistered = false
			AttentionLevel = 1
			QuizCount = 1
			if isCor {
				CorrectCount++
			}
			CorrectRate = CorrectCount / QuizCount
			_, err := s.client.RegisteredWord().
				Create().
				SetUserID(userID).
				SetWordID(x.WordID).
				SetIsActive(false).
				SetAttentionLevel(AttentionLevel).
				SetQuizCount(QuizCount).
				SetCorrectCount(CorrectCount).
				SetCorrectRate(CorrectRate).
				Save(ctx)
			if err != nil {
				return nil, errors.New("Failed to create RegisteredWord")
			}
		} else {
			CorrectCount = rw.CorrectCount
			QuizCount = rw.QuizCount
			if isCor {
				CorrectCount++
			}
			QuizCount++
			CorrectRate = CorrectCount / QuizCount
			_, err := rw.Update().
				SetQuizCount(QuizCount).
				SetCorrectCount(CorrectCount).
				SetCorrectRate(CorrectRate).
				Save(ctx)
			if err != nil {
				return nil, errors.New("Failed to update RegisteredWord")
			}
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
			ResisteredWord: models.ResisteredWord{
				IsRegistered:   IsRegistered,
				AttentionLevel: AttentionLevel,
				QuizCount:      QuizCount,
				CorrectCount:   CorrectCount,
			},
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
