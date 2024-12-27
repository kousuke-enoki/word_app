package exam_service

import (
	"context"
	"math/rand"
	"strconv"

	"word_app/backend/ent"
	"word_app/backend/ent/japanesemean"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/ent/wordinfo"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

func (s *ExamServiceImpl) CreateExam(ctx context.Context, req *models.CreateExamRequest) (*models.CreateExamResponse, error) {
	userID := req.UserID
	partsOfSpeeches := req.PartsOfSpeeches
	targetWords := req.TargetWords
	examCount := req.ExamCount

	// クエリのベースを作成
	wordQuery := s.client.Word.Query().
		WithWordInfos(func(wiQuery *ent.WordInfoQuery) {
			if len(partsOfSpeeches) > 0 {
				wiQuery.Where(wordinfo.PartOfSpeechIDIn(partsOfSpeeches...))
			}
			wiQuery.WithJapaneseMeans()
		}).
		WithRegisteredWords(func(rwQuery *ent.RegisteredWordQuery) {
			rwQuery.Where(registeredword.UserID(userID))
		})

	// 単語のターゲットが "registered" の場合
	if targetWords == "registered" {
		wordQuery = wordQuery.Where(word.HasRegisteredWordsWith(registeredword.UserID(userID)))
	}

	// ランダムに単語を取得
	words, err := wordQuery.Order(ent.Desc()).All(ctx) // 一旦全取得
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// シャッフルしてExamCount分だけ選択
	// rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })

	selectedWords := words
	if len(words) > examCount {
		selectedWords = words[:examCount]
	}

	// トランザクション開始
	tx, err := s.client.Tx(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	exam, err := tx.Exam.Create().
		SetUserID(userID).
		SetTotalQuestions(examCount).
		SetIsRunning(true).
		SetTargetWordTypes(targetWords).
		SetChoicesPosIds(partsOfSpeeches).
		Save(ctx)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return nil, err
	}

	var examQuestions []models.ExamQuestion
	for _, word := range selectedWords {
		// 正解の日本語訳
		if len(word.Edges.WordInfos) == 0 || len(word.Edges.WordInfos[0].Edges.JapaneseMeans) == 0 {
			continue // 日本語訳が存在しない場合はスキップ
		}

		correctMean := word.Edges.WordInfos[0].Edges.JapaneseMeans[0]
		correctJpmID := correctMean.ID

		// 誤回答（ランダムに3つ選択）
		allMeans, err := s.client.JapaneseMean.Query().Where(japanesemean.IDNEQ(correctJpmID)).All(ctx)
		if err != nil {
			logrus.Error(err)
			tx.Rollback()
			return nil, err
		}

		rand.Shuffle(len(allMeans), func(i, j int) { allMeans[i], allMeans[j] = allMeans[j], allMeans[i] })
		wrongMeans := allMeans
		if len(allMeans) > 3 {
			wrongMeans = allMeans[:3]
		}

		// 選択肢のIDをまとめる
		choiceIDs := []int{correctJpmID}
		for _, wm := range wrongMeans {
			choiceIDs = append(choiceIDs, wm.ID)
		}
		rand.Shuffle(len(choiceIDs), func(i, j int) { choiceIDs[i], choiceIDs[j] = choiceIDs[j], choiceIDs[i] })

		// ExamQuestionを作成
		examQuestion, err := tx.ExamQuestion.Create().
			SetExamID(exam.ID).
			SetCorrectJpmID(correctJpmID).
			SetChoicesJpmIds(choiceIDs).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// wrongMeans を models.QuestionJpm に変換
		var convertedWrongMeans []models.QuestionJpm
		for _, wm := range wrongMeans {
			convertedWrongMeans = append(convertedWrongMeans, models.QuestionJpm{
				JapaneseMeanID: wm.ID, // wm.ID を使用 (ent.JapaneseMean のフィールドに基づく)
				Name:           wm.Name,
			})
		}

		// examQuestions に追加
		examQuestions = append(examQuestions, models.ExamQuestion{
			ExamQuestionID: examQuestion.ID,
			WordName:       word.Name,
			QuestionJpms: append([]models.QuestionJpm{
				{JapaneseMeanID: correctMean.ID, Name: correctMean.Name},
			}, convertedWrongMeans...), // 変換済みスライスを追加
		})
	}

	// コミット
	if err := tx.Commit(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &models.CreateExamResponse{
		ExamID:        exam.ID,
		TotalExams:    strconv.Itoa(examCount),
		ExamQuestions: examQuestions,
	}, nil
}
