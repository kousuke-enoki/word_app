package quiz

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/quiz"
	"word_app/backend/ent/quizquestion"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

func (s *ServiceImpl) GetNextOrResume(
	ctx context.Context,
	userID int,
	req *models.GetQuizRequest,
) (*models.GetQuizResponse, error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		defer func() {
			if p := recover(); p != nil {
				logrus.Error(err)
				_ = tx.Rollback()
				panic(p)
			} else if err != nil {
				logrus.Error(err)
				_ = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}()
		return nil, err
	}
	var q *ent.Quiz

	q, err = tx.Quiz.
		Query().
		Where(quiz.UserID(userID), quiz.IsRunning(true)).
		Order(ent.Desc(quiz.FieldID)).
		First(ctx)
	if err != nil {
		return &models.GetQuizResponse{
			IsRunningQuiz: false,
		}, err
	}

	// questionNumber を決定
	var targetNum int
	if req.BeforeQuestionNumber != nil {
		targetNum = *req.BeforeQuestionNumber + 1
	} else {
		unanswered, qErr := tx.QuizQuestion.
			Query().
			Where(
				quizquestion.QuizIDEQ(q.ID),
				quizquestion.AnswerJpmIDIsNil(),
			).
			Order(ent.Asc(quizquestion.FieldQuestionNumber)).
			First(ctx)

		if qErr == nil {
			targetNum = unanswered.QuestionNumber
		} else if ent.IsNotFound(qErr) {
			targetNum = q.TotalQuestionsCount
		} else {
			return &models.GetQuizResponse{
				IsRunningQuiz: false,
			}, nil
		}
	}

	// 該当問題を取得
	qq, err := tx.QuizQuestion.
		Query().
		Where(
			quizquestion.QuizIDEQ(q.ID),
			quizquestion.QuestionNumber(targetNum),
		).
		Only(ctx)
	if err != nil {
		return &models.GetQuizResponse{
			IsRunningQuiz: false,
		}, nil
	}

	return &models.GetQuizResponse{
		IsRunningQuiz: true,
		NextQuestion: models.NextQuestion{
			QuizID:         q.ID,
			QuestionNumber: qq.QuestionNumber,
			WordName:       qq.WordName,
			ChoicesJpms:    qq.ChoicesJpms,
		},
	}, nil
}
