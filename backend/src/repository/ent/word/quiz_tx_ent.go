// src/repository/ent/quiz_tx_ent.go
package ent

import (
	"context"

	ent "word_app/backend/ent"
	rep "word_app/backend/src/repository"
)

func (r *wordRepoEnt) BeginTx(ctx context.Context) (rep.Tx, error) {
	tx, err := r.c.Tx(ctx)
	if err != nil {
		return nil, err
	}
	return &quizTxEnt{Tx: tx}, nil
}

type quizTxEnt struct {
	*ent.Tx // Ent のトランザクション
}

// -------------- CreateQuiz ----------------
func (q *quizTxEnt) CreateQuiz(ctx context.Context, userID int, rec rep.QuizRecord) (int, error) {
	id, err := q.Quiz.
		Create().
		SetUserID(userID).
		SetQuizNumber(0). // ここで採番しても良いし DB 側シーケンスでも可
		SetIsRunning(true).
		SetTotalQuestionsCount(rec.QuestionCount).
		SetIsSaveResult(rec.IsSaveResult).
		SetIsRegisteredWords(rec.IsRegisteredWords).
		SetSettingCorrectRate(rec.SettingCorrectRate).
		SetIsIdioms(rec.IsIdioms).
		SetIsSpecialCharacters(rec.IsSpecialChars).
		SetAttentionLevelList(rec.AttentionLevels).
		SetChoicesPosIds(rec.ChoicePosIDs).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return id.ID, nil
}

// -------------- CreateQuestions (Bulk) ----
func (q *quizTxEnt) CreateQuestions(ctx context.Context,
	quizID int, qs []rep.QuizQuestionRecord) error {

	bulk := make([]*ent.QuizQuestionCreate, len(qs))
	for i, r := range qs {
		bulk[i] = q.QuizQuestion.
			Create().
			SetQuizID(quizID).
			SetQuestionNumber(r.QuestionNumber).
			SetWordID(r.WordID).
			SetPosID(r.PosID).
			SetCorrectJpmID(r.CorrectJpmID).
			SetChoicesJpmsJSON(r.ChoicesJSON)
	}
	return q.Client().QuizQuestion.CreateBulk(bulk...).Exec(ctx)
}

// -------------- Commit / Rollback ----------
func (q *quizTxEnt) Commit() error   { return q.Tx.Commit() }
func (q *quizTxEnt) Rollback() error { return q.Tx.Rollback() }
