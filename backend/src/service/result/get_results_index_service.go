package result

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/quiz"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

func (s *ServiceImpl) GetSummaries(
	ctx context.Context,
	userID int,
) ([]models.ResultSummary, error) {

	rows, err := s.client.Quiz().
		Query().
		Where(
			quiz.UserID(userID),
			quiz.IsRunning(false),
			quiz.DeletedAtIsNil(),
		).
		Order(ent.Desc(quiz.FieldCreatedAt)). // 新しい順
		All(ctx)
	if err != nil {
		return nil, err
	}

	summaries := make([]models.ResultSummary, 0, len(rows))
	for _, q := range rows {
		summaries = append(summaries, models.ResultSummary{
			QuizNumber:          q.QuizNumber,
			CreatedAt:           q.CreatedAt,
			IsRegisteredWords:   q.IsRegisteredWords,
			IsIdioms:            q.IsIdioms,
			IsSpecialCharacters: q.IsSpecialCharacters,
			ChoicesPosIDs:       q.ChoicesPosIds,
			TotalQuestionsCount: q.TotalQuestionsCount,
			CorrectCount:        q.CorrectCount,
			ResultCorrectRate:   q.ResultCorrectRate,
		})
	}

	logrus.Debugf("results index: %+v", summaries)
	return summaries, nil
}
