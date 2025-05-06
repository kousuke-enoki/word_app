package quiz_service

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/quiz"
	"word_app/backend/ent/quizquestion"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

func (s *QuizServiceImpl) GetNextOrResume(
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

	if req.QuizID != nil { // クエリで指定あり
		q, err = tx.Quiz.
			Query().
			Where(quiz.IDEQ(*req.QuizID), quiz.UserID(userID)).
			Only(ctx)
	} else {
		q, err = tx.Quiz.
			Query().
			Where(quiz.UserID(userID), quiz.IsRunning(true)).
			Order(ent.Desc(quiz.FieldID)).
			First(ctx)
	}
	if err != nil {
		return &models.GetQuizResponse{
			IsRunningQuiz: false,
		}, err // not found → 呼び出し側で 204 などに
	}

	//----------------------------------------
	// 2. questionNumber を決定
	//----------------------------------------
	var targetNum int
	if req.BeforeQuestionNumber != nil {
		targetNum = *req.BeforeQuestionNumber + 1
	} else {
		unanswered, err := tx.QuizQuestion.
			Query().
			Where(
				quizquestion.QuizIDEQ(q.ID),
				quizquestion.AnswerJpmID(0),
			).
			Order(ent.Asc(quizquestion.FieldQuestionNumber)).
			First(ctx)

		if err == nil {
			targetNum = unanswered.QuestionNumber
		} else if ent.IsNotFound(err) {
			targetNum = q.TotalQuestionsCount
		} else {
			return &models.GetQuizResponse{
				IsRunningQuiz: false,
			}, nil
		}
	}

	//----------------------------------------
	// 3. 該当問題を取得
	//----------------------------------------
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

/*** helper ***/
func rollbackOrCommit(perr *error, tx *ent.Tx) {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if *perr != nil {
		_ = tx.Rollback()
	} else {
		*perr = tx.Commit()
	}
}
