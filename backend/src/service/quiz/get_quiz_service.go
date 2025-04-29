package quiz_service

import (
	"context"

	"word_app/backend/src/models"
)

func (s *QuizServiceImpl) GetQuiz(ctx context.Context, req *models.GetQuizRequest) (*models.GetQuizResponse, error) {
	// userID := req.UserID
	// quizID := req.QuizID

	// // // クエリのベースを作成
	// // wordQuery := s.client.Word.Query().
	// // 	WithWordInfos(func(wiQuery *ent.WordInfoQuery) {
	// // 		if len(partsOfSpeeches) > 0 {
	// // 			wiQuery.Where(wordinfo.PartOfSpeechIDIn(partsOfSpeeches...))
	// // 		}
	// // 		wiQuery.WithJapaneseMeans()
	// // 	}).
	// // 	WithRegisteredWords(func(rwQuery *ent.RegisteredWordQuery) {
	// // 		rwQuery.Where(registeredword.UserID(userID))
	// // 	})

	// // // 単語のターゲットが "all" の場合
	// // if targetWords == "all" {
	// // 	wordQuery = wordQuery.Where(word.HasRegisteredWordsWith(registeredword.UserID(userID)))
	// // }

	// // // ランダムに単語を取得
	// // words, err := wordQuery.Order(ent.Desc()).All(ctx) // 一旦全取得
	// // if err != nil {
	// // 	logrus.Error(err)
	// // 	return nil, err
	// // }

	// // // シャッフルしてQuizCount分だけ選択
	// // rand.Seed(time.Now().UnixNano())
	// // rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })

	// // selectedWords := words
	// // if len(words) > quizCount {
	// // 	selectedWords = words[:quizCount]
	// // }

	// // // トランザクション開始
	// // tx, err := s.client.Tx(ctx)
	// // if err != nil {
	// // 	logrus.Error(err)
	// // 	return nil, err
	// // }
	// // defer func() {
	// // 	if r := recover(); r != nil {
	// // 		tx.Rollback()
	// // 		panic(r)
	// // 	}
	// // }()
	// quiz, err := s.client.Quiz.
	// 	Query().
	// 	Where(quiz.ID(quizID)).
	// 	WithQuizQuestions(func(wq *ent.QuizQuestionQuery) {
	// 		wq.WithJapaneseMeans().WithPartOfSpeech()
	// 	}).
	// 	WithRegisteredWords().
	// 	Only(ctx)
	// quiz, err := s.client.Quiz.Get(ctx, quizID)
	// // SetUserID(userID).
	// // SetTotalQuestions(quizCount).
	// // SetIsRunning(true).
	// // SetTargetWordTypes(targetWords).
	// // SetChoicesPosIds(partsOfSpeeches).
	// // Save(ctx)
	// if err != nil {
	// 	return nil, errors.New("failed to fetch quiz")
	// }

	// var quizQuestions []models.QuizQuestion
	// for _, word := range selectedWords {
	// 	// 正解の日本語訳
	// 	if len(word.Edges.WordInfos) == 0 || len(word.Edges.WordInfos[0].Edges.JapaneseMeans) == 0 {
	// 		continue // 日本語訳が存在しない場合はスキップ
	// 	}

	// 	correctMean := word.Edges.WordInfos[0].Edges.JapaneseMeans[0]
	// 	correctJpmID := correctMean.ID

	// 	// 誤回答（ランダムに3つ選択）
	// 	allMeans, err := s.client.JapaneseMean.Query().Where(japanesemean.IDNEQ(correctJpmID)).All(ctx)
	// 	if err != nil {
	// 		logrus.Error(err)
	// 		tx.Rollback()
	// 		return nil, err
	// 	}

	// 	rand.Shuffle(len(allMeans), func(i, j int) { allMeans[i], allMeans[j] = allMeans[j], allMeans[i] })
	// 	wrongMeans := allMeans
	// 	if len(allMeans) > 3 {
	// 		wrongMeans = allMeans[:3]
	// 	}

	// 	// 選択肢のIDをまとめる
	// 	choiceIDs := []int{correctJpmID}
	// 	for _, wm := range wrongMeans {
	// 		choiceIDs = append(choiceIDs, wm.ID)
	// 	}
	// 	rand.Shuffle(len(choiceIDs), func(i, j int) { choiceIDs[i], choiceIDs[j] = choiceIDs[j], choiceIDs[i] })

	// 	// QuizQuestionを作成
	// 	quizQuestion, err := tx.QuizQuestion.Create().
	// 		SetQuizID(quiz.ID).
	// 		// SetWordID(word.ID).
	// 		SetAnswer(correctJpmID).
	// 		SetChoicesJpmIds(choiceIDs).
	// 		Save(ctx)
	// 	if err != nil {
	// 		logrus.Error(err)
	// 		tx.Rollback()
	// 		return nil, err
	// 	}

	// 	// wrongMeans を models.QuestionJpm に変換
	// 	var convertedWrongMeans []models.QuestionJpm
	// 	for _, wm := range wrongMeans {
	// 		convertedWrongMeans = append(convertedWrongMeans, models.QuestionJpm{
	// 			JapaneseMeanID: wm.ID, // wm.ID を使用 (ent.JapaneseMean のフィールドに基づく)
	// 			Name:           wm.Name,
	// 		})
	// 	}

	// 	// quizQuestions に追加
	// 	quizQuestions = append(quizQuestions, models.QuizQuestion{
	// 		QuizQuestionID: quizQuestion.ID,
	// 		WordName:       word.Name,
	// 		QuestionJpms: append([]models.QuestionJpm{
	// 			{JapaneseMeanID: correctMean.ID, Name: correctMean.Name},
	// 		}, convertedWrongMeans...), // 変換済みスライスを追加
	// 	})
	// }

	// // コミット
	// if err := tx.Commit(); err != nil {
	// 	logrus.Error(err)
	// 	return nil, err
	// }

	// return &models.CreateQuizResponse{
	// 	QuizID:        quiz.ID,
	// 	TotalQuizs:    strconv.Itoa(quizCount),
	// 	QuizQuestions: quizQuestions,
	// }, nil
	return nil, nil
}
