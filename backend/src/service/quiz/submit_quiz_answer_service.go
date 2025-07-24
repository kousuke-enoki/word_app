package quiz

import (
	"context"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/quiz"
	"word_app/backend/ent/quizquestion"
	"word_app/backend/src/models"
)

// SubmitAnswerAndRoute = 回答を保存し、次の分岐を判断
func (s *ServiceImpl) SubmitAnswerAndRoute(
	ctx context.Context,
	userID int,
	in *models.PostAnswerQuestionRequest,
) (res *models.AnswerRouteRes, err error) {
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
		UpdateOneID(qq.ID).
		SetAnswerJpmID(in.AnswerJpmID).
		SetIsCorrect(isCorrect).
		SetAnsweredAt(time.Now()).
		SetTimeMs(elapsedMs).
		Save(ctx); err != nil {
		return nil, err
	}

	// ③ 次問題があるか？
	nextQQ, errN := tx.QuizQuestion.
		Query().
		Where(
			quizquestion.QuizIDEQ(in.QuizID),
			quizquestion.QuestionNumber(in.QuestionNumber+1),
		).
		Only(ctx)
	if ent.IsNotFound(errN) {
		if finishedErr := s.finishQuizTx(ctx, tx, qq.Edges.Quiz, userID); finishedErr != nil {
			err = finishedErr
			return nil, err
		}
		res = &models.AnswerRouteRes{
			IsFinish:   true,
			IsCorrect:  isCorrect,
			QuizNumber: qq.Edges.Quiz.QuizNumber,
		}
		return
	} else if errN != nil {
		err = errN // defer に渡す
		return nil, err
	}
	// ---------- 次の問題あり ----------
	res = &models.AnswerRouteRes{
		IsFinish: false,
		NextQuestion: models.NextQuestion{
			QuizID:         in.QuizID,
			QuestionNumber: nextQQ.QuestionNumber,
			WordName:       nextQQ.WordName,
			ChoicesJpms:    nextQQ.ChoicesJpms,
		},
		IsCorrect: isCorrect,
	}
	return
}
