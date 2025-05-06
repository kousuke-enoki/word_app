package quiz_service

import (
	"context"
	"errors"
	"math/rand"

	"word_app/backend/ent"
	"word_app/backend/ent/japanesemean"
	"word_app/backend/ent/quiz"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	wordEnt "word_app/backend/ent/word"
	"word_app/backend/ent/wordinfo"
	"word_app/backend/src/models"

	"entgo.io/ent/dialect/sql"
	"github.com/sirupsen/logrus"
)

func (s *QuizServiceImpl) CreateQuiz(
	ctx context.Context,
	userID int,
	req *models.CreateQuizReq,
) (*models.CreateQuizResponse, error) {
	// トランザクション開始
	tx, err := s.client.Tx(ctx)
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

	questionCount := req.QuestionCount
	isSaveResult := req.IsSaveResult
	isRegisteredWords := req.IsRegisteredWords
	correctRate := req.CorrectRate
	attentionLevelList := req.AttentionLevelList
	partsOfSpeeches := req.PartsOfSpeeches
	isIdioms := req.IsIdioms
	isSpecialCharacters := req.IsSpecialCharacters

	quizCount, err := tx.Quiz.Query().Where(quiz.UserID(userID)).Count(ctx)
	if err != nil {
		return nil, err
	}

	quiz, err := tx.Quiz.
		Create().
		SetUserID(userID).
		SetQuizNumber(quizCount + 1).
		SetIsRunning(true).
		SetTotalQuestionsCount(questionCount).
		SetIsSaveResult(isSaveResult).
		SetIsRegisteredWords(isRegisteredWords).
		SetSettingCorrectRate(correctRate).
		SetIsIdioms(isIdioms).
		SetIsSpecialCharacters(isSpecialCharacters).
		SetAttentionLevelList(attentionLevelList).
		SetChoicesPosIds(partsOfSpeeches).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// クエリのベースを作成
	// wordQuery := tx.Word.Query().
	// 	WithWordInfos(func(wiQuery *ent.WordInfoQuery) {
	// 		if len(partsOfSpeeches) > 0 {
	// 			wiQuery.Where(wordinfo.PartOfSpeechIDIn(partsOfSpeeches...))
	// 		}
	// 		wiQuery.WithJapaneseMeans()
	// 	}).
	// 	WithRegisteredWords(func(rwQuery *ent.RegisteredWordQuery) {
	// 		rwQuery.Where(registeredword.UserID(userID))
	// 	})
	wordQuery := tx.Word.
		Query().
		WithWordInfos(func(wi *ent.WordInfoQuery) {
			if len(partsOfSpeeches) > 0 {
				wi.Where(wordinfo.PartOfSpeechIDIn(partsOfSpeeches...))
			}
			wi.WithJapaneseMeans()
		}).
		WithRegisteredWords(func(rw *ent.RegisteredWordQuery) {
			rw.Where(registeredword.UserID(userID))
		})

		// 登録単語フィルタ
	switch isRegisteredWords {
	case 1: // registered only
		wordQuery = wordQuery.
			Where(word.HasRegisteredWordsWith(
				registeredword.UserID(userID),
				registeredword.IsActiveEQ(true),
				registeredword.CorrectCountLTE(correctRate),
				registeredword.AttentionLevelIn(attentionLevelList...),
			))
	case 2: // unregistered only
		wordQuery = wordQuery.
			Where(word.Not(word.HasRegisteredWordsWith(
				registeredword.UserID(userID),
				registeredword.IsActiveEQ(true),
			)))
	}
	// 慣用句／特殊単語フィルタ
	if isIdioms != 0 {
		wordQuery = wordQuery.Where(word.IsIdiomsEQ(isIdioms == 1))
	}
	if isSpecialCharacters != 0 {
		wordQuery = wordQuery.Where(word.IsSpecialCharactersEQ(isSpecialCharacters == 1))
	}
	// words, err := wordQuery.
	// Order(ent.OrderFunc("RANDOM()")). // ORDER BY RANDOM()
	// Limit(questionCount).             // LIMIT N
	// All(ctx)
	words, err := wordQuery.
		Order(func(s *sql.Selector) {
			s.OrderBy("RANDOM()")
		}).                   // ORDER BY RANDOM()
		Limit(questionCount). // LIMIT N
		All(ctx)
	if err != nil {
		return nil, err
	}
	if len(words) < questionCount {
		return nil, errors.New("quiz question is not enough")
	}

	// 作成できたQuizquestionの数
	createdQuestionCount := 0
	// 一問目用のres
	firstQuizQuestions := models.NextQuestion{QuestionNumber: 1}
	for index, word := range words {
		logrus.Debugf("word %+v", word)
		// 正解の日本語訳
		if len(word.Edges.WordInfos) == 0 || len(word.Edges.WordInfos[0].Edges.JapaneseMeans) == 0 {
			continue // 日本語訳が存在しない場合はスキップ
		}

		correctMean := word.Edges.WordInfos[0].Edges.JapaneseMeans[0]
		correctMeanPosID := word.Edges.WordInfos[0].PartOfSpeechID
		correctJpmID := correctMean.ID

		// 誤回答（ランダムに3つ選択）
		choicesJpmQuery := tx.JapaneseMean.
			Query().
			Where(japanesemean.IDNEQ(correctJpmID)).
			// 品詞一致
			Where(japanesemean.HasWordInfoWith(
				wordinfo.PartOfSpeechID(correctMeanPosID),
			))

		if isIdioms != 0 {
			choicesJpmQuery = choicesJpmQuery.
				Where(japanesemean.HasWordInfoWith(
					wordinfo.HasWordWith(wordEnt.IsIdioms(isIdioms == 1)),
				))
		}
		if isSpecialCharacters != 0 {
			choicesJpmQuery = choicesJpmQuery.
				Where(japanesemean.HasWordInfoWith(
					wordinfo.HasWordWith(wordEnt.IsSpecialCharactersEQ(isSpecialCharacters == 1)),
				))
		}
		// // 慣用句を条件にする場合
		// choicesJpmQuery = choicesJpmQuery.Where(func(s *sql.Selector) {
		// 	s.Where(sql.P(func(b *sql.Builder) {
		// 		b.WriteString("EXISTS (SELECT 1 FROM words WHERE words.id = japanese_means.word_id AND words.is_idioms = ?)")
		// 		b.Arg(word.IsIdioms)
		// 	}))
		// })
		// // 特殊単語を条件にする場合
		// choicesJpmQuery = choicesJpmQuery.Where(func(s *sql.Selector) {
		// 	s.Where(sql.P(func(b *sql.Builder) {
		// 		b.WriteString("EXISTS (SELECT 1 FROM words WHERE words.id = japanese_means.word_id AND words.is_special_characters = ?)")
		// 		b.Arg(word.IsSpecialCharacters)
		// 	}))
		// })
		// ⭐️ 乱数 3 件だけ取得
		// wrongMeans, err := choicesJpmQuery.
		// 	Order(ent.OrderFunc("RANDOM()")). // ORDER BY RANDOM()
		// 	Limit(3).
		// 	All(ctx)
		// if err != nil {
		// 	return nil, err
		wrongMeans, err := choicesJpmQuery.
			Order(func(s *sql.Selector) {
				s.OrderBy("RANDOM()")
			}).
			Limit(3).
			All(ctx)
		if err != nil {
			return nil, err
		}
		// allMeans, err := choicesJpmQuery.All(ctx)
		// if err != nil {
		// 	logrus.Error(err)
		// 	return nil, err
		// }

		// // wrongMeansは間違いの選択肢
		// rand.Shuffle(len(allMeans), func(i, j int) { allMeans[i], allMeans[j] = allMeans[j], allMeans[i] })
		// wrongMeans := allMeans
		// if len(allMeans) > 3 {
		// 	wrongMeans = allMeans[:3]
		// }

		// 選択肢のchoicesJpmをまとめる
		// choicesJpmsは選択肢のjpmIDとjpm.nameのオブジェクトリスト,firstQuizQuestionのフロント返却用
		var choicesJpms []models.ChoiceJpm
		// 選択肢のjpm.idをまとめる
		// choicesJpmIdsは選択肢のjpmIDリスト,テーブル保存用
		// var choicesJpmIds []int
		// index==0(一問目の場合)フロントに返却用の問題をresに入れる処理を通す
		// if index == 0 {
		for _, wm := range wrongMeans {
			choicesJpms = append(choicesJpms, models.ChoiceJpm{
				JapaneseMeanID: wm.ID,
				Name:           wm.Name,
			})
		}
		var correctJpm = models.ChoiceJpm{
			JapaneseMeanID: correctMean.ID,
			Name:           correctMean.Name,
		}

		choicesJpms = append(choicesJpms, correctJpm)
		rand.Shuffle(len(choicesJpms), func(i, j int) { choicesJpms[i], choicesJpms[j] = choicesJpms[j], choicesJpms[i] })

		// for _, wm := range choicesJpms {
		// 	choicesJpmIds = append(choicesJpmIds, wm.JapaneseMeanID)
		// }
		// } else {
		// 	for _, wm := range wrongMeans {
		// 		choicesJpmIds = append(choicesJpmIds, wm.ID)
		// 	}

		// 	choicesJpmIds = append(choicesJpmIds, correctJpmID)
		// 	rand.Shuffle(len(choicesJpmIds), func(i, j int) { choicesJpmIds[i], choicesJpmIds[j] = choicesJpmIds[j], choicesJpmIds[i] })

		// }

		// QuizQuestionを作成
		_, err = tx.QuizQuestion.Create().
			SetQuizID(quiz.ID).
			SetQuestionNumber(index + 1).
			SetWordID(word.ID).
			SetWordName(word.Name).
			SetPosID(correctMeanPosID).
			SetCorrectJpmID(correctJpmID).
			SetChoicesJpms(choicesJpms).
			Save(ctx)
		if err != nil {
			return nil, err
		} else {
			createdQuestionCount++
		}

		firstQuizQuestions = models.NextQuestion{
			QuestionNumber: 1,
			WordName:       word.Name,
			ChoicesJpms:    choicesJpms,
		}
		// var convertedWrongMeans []models.ChoicesJpm
		// for _, wm := range wrongMeans {
		// 	convertedWrongMeans = append(convertedWrongMeans, models.ChoicesJpm{
		// 		JapaneseMeanID: wm.ID, // wm.ID を使用 (ent.JapaneseMean のフィールドに基づく)
		// 		Name:           wm.Name,
		// 	})
		// }

		// // quizQuestions
		// quizQuestions = append(quizQuestions, models.QuizQuestion{
		// 	QuizID:         quizQuestion.ID,
		// 	QuestionNumber: quizQuestion.QuestionNumber,
		// 	WordName:       word.Name,
		// 	ChoicesJpms:    choicesJpms,
		// })

	}
	if createdQuestionCount < req.QuestionCount {
		return nil, errors.New("quiz question is not enough")
	}
	// コミット
	if err := tx.Commit(); err != nil {
		logrus.Error(err)
		return nil, err
	}
	// firstQuizQuestion, err := s.client.QuizQuestion().
	// 	Query().
	// 	Where(quizquestion.QuizID(quiz.ID), quizquestion.QuestionNumber(1)).
	// 	WithWords().
	// 	With().
	// 	Only(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// // 選択肢のchoicesJpmをまとめる
	// // choicesJpmsは選択肢のjpmIDとjpm.nameのオブジェクトリスト,フロント返却用
	// var choicesJpms []models.ChoicesJpm
	// for _, jpmID := range firstQuizQuestion.ChoicesJpmIDs {

	// 	choicesJpms = append(choicesJpms, models.ChoicesJpm{
	// 		JapaneseMeanID: wm.ID,
	// 		Name:           wm.Name,
	// 	})
	// }
	// var correctJpm = models.ChoicesJpm{
	// 	JapaneseMeanID: correctMean.ID,
	// 	Name:           correctMean.Name,
	// }

	// QuizQuestionRes := models.QuizQuestionRes{
	// 	QuestionNumber: 1,
	// 	WordName: firstQuizQuestion.edge.word.name,
	// 	ChoicesJpms: ,
	// }

	return &models.CreateQuizResponse{
		QuizID:               quiz.ID,
		TotalCreatedQuestion: createdQuestionCount,
		NextQuestion:         firstQuizQuestions,
	}, nil
}
