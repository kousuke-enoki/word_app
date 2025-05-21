package quiz_service

import (
	"context"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/quiz"
	"word_app/backend/ent/quizquestion"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

// SubmitAnswerAndRoute = 回答を保存し、次の分岐を判断
func (s *QuizServiceImpl) SubmitAnswerAndRoute(
	ctx context.Context,
	userID int,
	in *models.PostAnswerQuestionRequest,
) (res *models.AnswerRouteRes, err error) {
	logrus.Info("submitservice")
	logrus.Info(userID)
	logrus.Info(in)
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer finishTx(&err, tx)
	logrus.Info("1")
	logrus.Info(in.QuizID)
	logrus.Info(in.QuestionNumber)
	logrus.Info(quiz.UserID(userID))
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
	logrus.Info(qq)

	isCorrect := qq.CorrectJpmID == in.AnswerJpmID
	elapsedMs := int(time.Since(qq.CreatedAt).Milliseconds())

	logrus.Info(isCorrect)
	logrus.Info(in.AnswerJpmID)
	logrus.Info("2")
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

	logrus.Info("3")
	// ③ 次問題があるか？
	if nextQQ, errN := tx.QuizQuestion.
		Query().
		Where(
			quizquestion.QuizIDEQ(in.QuizID),
			quizquestion.QuestionNumber(in.QuestionNumber+1),
		).
		Only(ctx); errN == nil { // ★ errN に変更
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
		return // err は defer で上書き
	} else if ent.IsNotFound(errN) { // ---- これが最終 ----
		if finishedErr := s.finishQuizTx(ctx, tx, qq.Edges.Quiz, userID); finishedErr != nil {
			err = finishedErr // defer に渡す
			return nil, err
		}

		res = &models.AnswerRouteRes{
			IsFinish:   true,
			IsCorrect:  isCorrect,
			QuizNumber: qq.Edges.Quiz.QuizNumber,
		}
		return // Result は返さない

	} else {
		err = errN // ★ defer に渡す
		return nil, err
	}
}
