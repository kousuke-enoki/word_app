package result_service

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/quiz"
	"word_app/backend/ent/quizquestion"
	"word_app/backend/ent/registeredword"
	"word_app/backend/src/models"
)

/*
----------------------------------------------------------

	GET /results/:quizNo
	--------------------------------------------------------
*/
func (s *ResultServiceImpl) GetResultByQuizNo(
	ctx context.Context,
	userID int,
	quizNo int,
) (*models.Result, error) {

	q, err := s.client.
		Quiz().
		Query().
		Where(
			quiz.UserID(userID),
			quiz.QuizNumber(quizNo),
			quiz.IsRunning(false),
		).
		WithQuizQuestions(func(q *ent.QuizQuestionQuery) {
			q.Order(ent.Asc(quizquestion.FieldQuestionNumber))
		}).
		Only(ctx)
	if err != nil {
		return nil, err // not-found → ハンドラで 404
	}

	/* ---------- Question → DTO ---------- */
	resultQs := make([]models.ResultQuestion, 0, len(q.Edges.QuizQuestions))
	correctCnt := 0

	for _, qq := range q.Edges.QuizQuestions {
		isCor := qq.IsCorrect != nil && *qq.IsCorrect
		if isCor {
			correctCnt++
		}

		/* RegisteredWord (集計用情報) */
		rw, _ := s.client.RegisteredWord().
			Query().
			Where(
				registeredword.UserID(userID),
				registeredword.WordID(qq.WordID),
			).
			Only(ctx) // not-found → nil

		resQ := models.ResultQuestion{
			QuestionNumber: qq.QuestionNumber,
			WordID:         qq.WordID,
			WordName:       qq.WordName,
			PosID:          qq.PosID,
			CorrectJpmId:   qq.CorrectJpmID,
			ChoicesJpms:    qq.ChoicesJpms,
			AnswerJpmId:    derefInt(qq.AnswerJpmID),
			IsCorrect:      isCor,
			TimeMs:         derefInt(qq.TimeMs),
			RegisteredWord: models.RegisteredWord{
				IsRegistered:   rw != nil && rw.IsActive,
				AttentionLevel: ifNil(rw, 1, rw.AttentionLevel),
				QuizCount:      ifNil(rw, 0, rw.QuizCount),
				CorrectCount:   ifNil(rw, 0, rw.CorrectCount),
			},
		}
		resultQs = append(resultQs, resQ)
	}

	return &models.Result{
		QuizNumber:          q.QuizNumber,
		TotalQuestionsCount: q.TotalQuestionsCount,
		CorrectCount:        correctCnt,
		ResultCorrectRate:   percent(correctCnt, q.TotalQuestionsCount),
		ResultSetting: models.ResultSetting{
			IsSaveResult:        q.IsSaveResult,
			IsRegisteredWords:   q.IsRegisteredWords,
			SettingCorrectRate:  q.SettingCorrectRate,
			IsIdioms:            q.IsIdioms,
			IsSpecialCharacters: q.IsSpecialCharacters,
			AttentionLevelList:  q.AttentionLevelList,
			ChoicesPosIds:       q.ChoicesPosIds,
		},
		ResultQuestions: resultQs,
	}, nil
}

/* ---------- small helpers ---------- */

// func derefInt(ptr *int) int {
// 	if ptr == nil {
// 		return 0
// 	}
// 	return *ptr
// }

func ifNil[T any](rw *ent.RegisteredWord, def T, v T) T {
	if rw == nil {
		return def
	}
	return v
}

func percent(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) * 100 / float64(b)
}

// ヘルパ
func derefInt(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}
