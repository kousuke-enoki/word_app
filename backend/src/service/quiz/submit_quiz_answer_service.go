package quiz_service

import (
	"context"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/quiz"
	"word_app/backend/ent/quizquestion"
	"word_app/backend/src/models"
)

// SubmitAnswerAndRoute = 回答を保存し、次の分岐を判断
func (s *QuizServiceImpl) SubmitAnswerAndRoute(
	ctx context.Context,
	userID int,
	in *models.PostAnswerQuestionRequest,
) (*models.AnswerRouteRes, error) {

	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer finishTx(&err, tx)

	// ① 該当問題を取得
	qq, err := tx.QuizQuestion.
		Query().
		Where(
			quizquestion.QuizIDEQ(in.QuizID),
			quizquestion.QuestionNumber(in.QuestionNumber),
			quizquestion.HasQuizWith(quiz.UserID(userID)),
		).
		WithQuiz().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	isCorrect := qq.CorrectJpmID == in.AnswerJpmID
	elapsedMs := int(time.Since(qq.CreatedAt).Milliseconds())

	// ② 回答を保存
	if _, err = tx.QuizQuestion.
		UpdateOne(qq).
		SetAnswerJpmID(in.AnswerJpmID).
		SetIsCorrect(isCorrect).
		SetAnsweredAt(time.Now()).
		SetTimeMs(elapsedMs).
		Save(ctx); err != nil {
		return nil, err
	}

	// ③ 次問題があるか？
	nextQQ, err := tx.QuizQuestion.
		Query().
		Where(
			quizquestion.QuizIDEQ(in.QuizID),
			quizquestion.QuestionNumber(in.QuestionNumber+1),
		).
		Only(ctx)

	switch {
	case err == nil: // ------ 次の問題あり ------
		return &models.AnswerRouteRes{
			IsFinish: false,
			NextQuestion: models.NextQuestion{
				QuestionNumber: nextQQ.QuestionNumber,
				WordName:       nextQQ.WordName,
				ChoicesJpms:    nextQQ.ChoicesJpms,
			},
			IsCorrect: isCorrect,
		}, nil

	case ent.IsNotFound(err): // ------ これが最終 ------
		// FinishQuiz を呼び出して結果 DTO を構築
		result, err2 := s.finishQuizTx(ctx, tx, qq.Edges.Quiz)
		if err2 != nil {
			return nil, err2
		}

		return &models.AnswerRouteRes{
			IsFinish:  true,
			Result:    *result,
			IsCorrect: isCorrect,
		}, nil

	default:
		return nil, err
	}
}
