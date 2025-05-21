package quiz

import (
	"context"
	"math/rand"
	"word_app/backend/src/models"
)

type CreateQuizUsecase struct {
	Repo repository.WordRepo
	Rand *rand.Rand
}

func (u *CreateQuizUsecase) Execute(ctx context.Context, userID int, in *models.CreateQuizDTO) (*models.CreateQuizResponse, error) {

	// 1. 候補単語取得
	words, err := u.Repo.RandomSelectableWords(ctx, userID, buildFilter(in), in.QuestionCount)
	if err != nil {
		return nil, err
	}
	if len(words) < in.QuestionCount {
		return nil, domain.ErrNotEnough
	}

	// 2. Tx 開始
	tx, _ := u.Repo.BeginTx(ctx)
	defer tx.Rollback() // safety

	quizID, _ := tx.CreateQuiz(ctx, userID, toQuizRecord(in))
	questions := make([]QuizQuestionRecord, 0, len(words))
	// … slices.Map で変換 …
	if err := tx.CreateQuestions(ctx, quizID, questions); err != nil {
		return nil, err
	}
	_ = tx.Commit()

	// 3. DTO 組み立て
	return &models.CreateQuizResponse{
		QuizID: quizID, TotalCreatedQuestion: in.QuestionCount,
		NextQuestion: toNextQuestionDTO(questions[0]),
	}, nil
}
