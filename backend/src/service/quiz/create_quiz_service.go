package quiz_service

import (
	"context"
	"errors"
	"fmt"

	"math/rand"

	"word_app/backend/ent"
	"word_app/backend/ent/japanesemean"
	"word_app/backend/ent/quiz"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/ent/wordinfo"
	"word_app/backend/src/models"

	"entgo.io/ent/dialect/sql"
)

/*==================== public service ====================*/

func (s *QuizServiceImpl) CreateQuiz(
	ctx context.Context,
	userID int,
	req *models.CreateQuizReq,
) (resp *models.CreateQuizResponse, err error) {

	// ① ドメイン集約 ― 候補単語
	words, err := s.fetchCandidates(ctx, userID, req)
	if err != nil {
		return
	}

	// ② Tx で Quiz 作成ユースケース実行
	err = s.withTx(ctx, func(tx *ent.Tx) error {
		qEnt, err := s.ensureQuizRecord(ctx, tx, userID, req)
		if err != nil {
			return err
		}

		first, err := s.generateQuestions(ctx, tx, qEnt, words)
		if err != nil {
			return err
		}

		resp = &models.CreateQuizResponse{
			QuizID:               qEnt.ID,
			TotalCreatedQuestion: req.QuestionCount,
			NextQuestion:         first,
		}
		return nil
	})
	return
}

/*==================== “repository-like” helpers ====================*/

// fetchCandidates は「クエリ組み立て＋件数チェック」だけ担当。
// 将来 WordRepository に移設すれば service 層は一切修正不要。
func (s *QuizServiceImpl) fetchCandidates(
	ctx context.Context,
	userID int,
	req *models.CreateQuizReq,
) ([]*ent.Word, error) {

	q := s.baseWordQuery(userID, req) // クエリビルダーも分離
	words, err := q.
		Order(func(s *sql.Selector) { s.OrderBy("RANDOM()") }).
		Limit(req.QuestionCount).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("get words: %w", err)
	}
	if len(words) < req.QuestionCount {
		return nil, errors.New("quiz question is not enough")
	}
	return words, nil
}

// ensureQuizRecord は「排他チェック＋Quiz行作成」をまとめた 1 ユースケース。
// ここも将来 QuizRepository へ切り出し可能。
func (s *QuizServiceImpl) ensureQuizRecord(
	ctx context.Context,
	tx *ent.Tx,
	userID int,
	req *models.CreateQuizReq,
) (*ent.Quiz, error) {

	exists, err := tx.Quiz.Query().
		Where(quiz.UserID(userID), quiz.IsRunning(true)).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("another quiz is running: userID=%d", userID)
	}

	qCount, err := tx.Quiz.Query().Where(quiz.UserID(userID)).Count(ctx)
	if err != nil {
		return nil, err
	}

	return tx.Quiz.Create().
		SetUserID(userID).
		SetQuizNumber(qCount + 1).
		SetIsRunning(true).
		SetTotalQuestionsCount(req.QuestionCount).
		SetIsSaveResult(req.IsSaveResult).
		SetIsRegisteredWords(req.IsRegisteredWords).
		SetSettingCorrectRate(req.CorrectRate).
		SetIsIdioms(req.IsIdioms).
		SetIsSpecialCharacters(req.IsSpecialCharacters).
		SetAttentionLevelList(req.AttentionLevelList).
		SetChoicesPosIds(req.PartsOfSpeeches).
		Save(ctx)
}

// generateQuestions は「ドメインロジック：誤答抽出＋問題行作成」
// 今は Ent を直叩きだが、あとで QuestionRepository に。
func (s *QuizServiceImpl) generateQuestions(
	ctx context.Context,
	tx *ent.Tx,
	qEnt *ent.Quiz,
	words []*ent.Word,
) (models.NextQuestion, error) {
	var first models.NextQuestion

	for i, w := range words {
		wi := w.Edges.WordInfos[0]
		correct := wi.Edges.JapaneseMeans[0]

		wrongs, err := tx.JapaneseMean.Query().
			Where(
				japanesemean.IDNEQ(correct.ID),
				japanesemean.HasWordInfoWith(wordinfo.PartOfSpeechID(wi.PartOfSpeechID)),
			).
			Order(func(s *sql.Selector) { s.OrderBy("RANDOM()") }).
			Limit(3).
			All(ctx)
		if err != nil {
			return first, err
		}

		choices := buildChoices(correct, wrongs)

		qq, err := tx.QuizQuestion.Create().
			SetQuizID(qEnt.ID).
			SetQuestionNumber(i + 1).
			SetWordID(w.ID).
			SetWordName(w.Name).
			SetPosID(wi.PartOfSpeechID).
			SetCorrectJpmID(correct.ID).
			SetChoicesJpms(choices).
			Save(ctx)
		if err != nil {
			return first, err
		}

		if i == 0 {
			first = models.NextQuestion{
				QuizID:         qEnt.ID,
				QuestionNumber: qq.QuestionNumber,
				WordName:       qq.WordName,
				ChoicesJpms:    qq.ChoicesJpms,
			}
		}
	}
	return first, nil
}

/*==================== “query builder” ====================*/

// baseWordQuery はフィルタ条件を組み立てるだけ。
// 呼び出し側が order/limit を決められる＝再利用しやすい。
func (s *QuizServiceImpl) baseWordQuery(
	userID int,
	req *models.CreateQuizReq,
) *ent.WordQuery {

	q := s.client.Word().
		Query().
		Where(word.HasWordInfosWith(
			wordinfo.PartOfSpeechIDIn(req.PartsOfSpeeches...),
			wordinfo.HasJapaneseMeans(),
		)).
		WithWordInfos(func(wi *ent.WordInfoQuery) {
			wi.Where(wordinfo.PartOfSpeechIDIn(req.PartsOfSpeeches...)).
				WithJapaneseMeans()
		}).
		WithRegisteredWords(func(rw *ent.RegisteredWordQuery) {
			rw.Where(registeredword.UserID(userID))
		})

	switch req.IsRegisteredWords {
	case 1:
		q = q.Where(word.HasRegisteredWordsWith(
			registeredword.UserID(userID),
			registeredword.IsActiveEQ(true),
			registeredword.CorrectRateLTE(req.CorrectRate),
			registeredword.AttentionLevelIn(req.AttentionLevelList...),
		))
	case 2:
		q = q.Where(word.Not(word.HasRegisteredWordsWith(
			registeredword.UserID(userID),
			registeredword.IsActiveEQ(true))))
	}
	if flag := req.IsIdioms; flag != 0 {
		q = q.Where(word.IsIdiomsEQ(flag == 1))
	}
	if flag := req.IsSpecialCharacters; flag != 0 {
		q = q.Where(word.IsSpecialCharactersEQ(flag == 1))
	}
	return q
}

/*==================== tx wrapper & utility ====================*/

func (s *QuizServiceImpl) withTx(
	ctx context.Context,
	fn func(tx *ent.Tx) error,
) (err error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = fn(tx)
	return
}

func buildChoices(correct *ent.JapaneseMean, wrongs []*ent.JapaneseMean) []models.ChoiceJpm {
	choices := make([]models.ChoiceJpm, 0, 4)
	for _, jm := range wrongs {
		choices = append(choices, models.ChoiceJpm{JapaneseMeanID: jm.ID, Name: jm.Name})
	}
	choices = append(choices, models.ChoiceJpm{JapaneseMeanID: correct.ID, Name: correct.Name})
	rand.Shuffle(len(choices), func(i, j int) { choices[i], choices[j] = choices[j], choices[i] })
	return choices
}
