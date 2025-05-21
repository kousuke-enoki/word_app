package quiz_service

import (
	"context"
	"math/rand"

	"word_app/backend/src/models"
)

type CreateQuizUsecase struct {
	Repo WordRepo
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

// func (s *QuizServiceImpl) CreateQuiz(
// 	ctx context.Context,
// 	userID int,
// 	req *models.CreateQuizReq,
// ) (resp *models.CreateQuizResponse, err error) { // ★ named return
// 	logrus.Info(req)
// 	logrus.Info(req.PartsOfSpeeches)
// 	/* -------------------------------------------------
// 	① 先に出題候補を取得して “問題不足” を弾く
// 	-------------------------------------------------*/
// 	baseWQ := s.client.Word().
// 		Query().
// 		Where(
// 			word.HasWordInfosWith(
// 				wordinfo.PartOfSpeechIDIn(req.PartsOfSpeeches...),
// 				wordinfo.HasJapaneseMeans(), // ← wordInfo, jpmがないものを含むとエラーになるので除外
// 			),
// 		).
// 		WithWordInfos(func(wi *ent.WordInfoQuery) {
// 			wi.Where(wordinfo.PartOfSpeechIDIn(req.PartsOfSpeeches...)).
// 				WithJapaneseMeans()
// 		}).
// 		WithRegisteredWords(func(rw *ent.RegisteredWordQuery) {
// 			rw.Where(registeredword.UserID(userID))
// 		})
// 	logrus.Info("1")
// 	logrus.Info(baseWQ)
// 	// フラグ系フィルタ
// 	logrus.Info(req.CorrectRate)
// 	logrus.Info(req.AttentionLevelList)
// 	if req.IsRegisteredWords == 1 {
// 		baseWQ = baseWQ.Where(
// 			word.HasRegisteredWordsWith(
// 				registeredword.UserID(userID),
// 				registeredword.IsActiveEQ(true),
// 				registeredword.CorrectRateLTE(req.CorrectRate),
// 				registeredword.AttentionLevelIn(req.AttentionLevelList...),
// 			))
// 	}
// 	logrus.Info("2")
// 	if req.IsRegisteredWords == 2 {
// 		baseWQ = baseWQ.
// 			Where(word.Not(word.HasRegisteredWordsWith(
// 				registeredword.UserID(userID),
// 				registeredword.IsActiveEQ(true))))
// 	}
// 	logrus.Info("3")
// 	if req.IsIdioms != 0 {
// 		baseWQ = baseWQ.Where(word.IsIdiomsEQ(req.IsIdioms == 1))
// 	}
// 	logrus.Info("4")
// 	if req.IsSpecialCharacters != 0 {
// 		baseWQ = baseWQ.Where(word.IsSpecialCharactersEQ(req.IsSpecialCharacters == 1))
// 	}
// 	logrus.Info("5")

// 	// ★ まだトランザクションを張らずに取得
// 	words, err := baseWQ.
// 		Order(func(s *sql.Selector) { s.OrderBy("RANDOM()") }).
// 		Limit(req.QuestionCount).
// 		All(ctx)
// 	logrus.Info(words)
// 	if err != nil {
// 		logrus.Error(err)
// 		err = fmt.Errorf("get words: %w", err)
// 		return
// 	}
// 	if len(words) < req.QuestionCount { // ★ 問題数が足りなければ直帰
// 		err = errors.New("quiz question is not enough")
// 		return
// 	}

// 	logrus.Info("6")
// 	/* -------------------------------------------------
// 	② ここからトランザクション開始
// 	-------------------------------------------------*/
// 	tx, err := s.client.Tx(ctx)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}
// 	defer func() { // ★ named-err を見る
// 		if p := recover(); p != nil {
// 			_ = tx.Rollback()
// 			panic(p)
// 		}
// 		if err != nil {
// 			_ = tx.Rollback()
// 		} else {
// 			err = tx.Commit()
// 		}
// 	}()

// 	// --- 同時実行中のクイズチェック --------------------
// 	exists, err := tx.Quiz.Query().
// 		Where(quiz.UserID(userID), quiz.IsRunning(true)).
// 		Exist(ctx)
// 	logrus.Info(exists)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}
// 	if exists {
// 		err = fmt.Errorf("another quiz is running: userID=%d", userID)
// 		return
// 	}

// 	// --- quiz レコード INSERT -------------------------
// 	qCount, err := tx.Quiz.Query().Where(quiz.UserID(userID)).Count(ctx)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}
// 	logrus.Info(qCount)

// 	qEnt, err := tx.Quiz.
// 		Create().
// 		SetUserID(userID).
// 		SetQuizNumber(qCount + 1).
// 		SetIsRunning(true).
// 		SetTotalQuestionsCount(req.QuestionCount).
// 		SetIsSaveResult(req.IsSaveResult).
// 		SetIsRegisteredWords(req.IsRegisteredWords).
// 		SetSettingCorrectRate(req.CorrectRate).
// 		SetIsIdioms(req.IsIdioms).
// 		SetIsSpecialCharacters(req.IsSpecialCharacters).
// 		SetAttentionLevelList(req.AttentionLevelList).
// 		SetChoicesPosIds(req.PartsOfSpeeches).
// 		Save(ctx)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}
// 	logrus.Info(qEnt)

// 	/* -------------------------------------------------
// 	③ quiz_question を作成
// 	-------------------------------------------------*/
// 	var firstDTO models.NextQuestion
// 	for i, w := range words {
// 		if len(w.Edges.WordInfos) == 0 ||
// 			len(w.Edges.WordInfos[0].Edges.JapaneseMeans) == 0 {
// 			err = fmt.Errorf("invalid word (%d) selected – no meanings", w.ID)
// 			return
// 		}
// 		logrus.Info(w)
// 		wi := w.Edges.WordInfos[0] // 前段で WithWordInfos してある
// 		correct := wi.Edges.JapaneseMeans[0]

// 		// --- 誤答３件抽出 (tx を使う) ---------------
// 		wrongs, err2 := tx.JapaneseMean.Query().
// 			Where(
// 				japanesemean.IDNEQ(correct.ID),
// 				japanesemean.HasWordInfoWith(wordinfo.PartOfSpeechID(wi.PartOfSpeechID)),
// 			).
// 			Order(func(s *sql.Selector) { s.OrderBy("RANDOM()") }).
// 			Limit(3).
// 			All(ctx)
// 		if err2 != nil {
// 			logrus.Error(err2)
// 			err = err2
// 			return
// 		}
// 		logrus.Info(wrongs)

// 		// --- choices 組立 -----------------------------
// 		choices := make([]models.ChoiceJpm, 0, 4)
// 		for _, jm := range wrongs {
// 			logrus.Info(jm)
// 			choices = append(choices, models.ChoiceJpm{JapaneseMeanID: jm.ID, Name: jm.Name})
// 		}
// 		choices = append(choices, models.ChoiceJpm{JapaneseMeanID: correct.ID, Name: correct.Name})
// 		logrus.Info(choices)
// 		rand.Shuffle(len(choices), func(i, j int) { choices[i], choices[j] = choices[j], choices[i] })

// 		// --- INSERT quiz_question ----------------------
// 		qq, err2 := tx.QuizQuestion.Create().
// 			SetQuizID(qEnt.ID).
// 			SetQuestionNumber(i + 1).
// 			SetWordID(w.ID).
// 			SetWordName(w.Name).
// 			SetPosID(wi.PartOfSpeechID).
// 			SetCorrectJpmID(correct.ID).
// 			SetChoicesJpms(choices).
// 			Save(ctx)
// 		if err2 != nil {
// 			logrus.Error(err2)
// 			err = err2
// 			return
// 		}
// 		logrus.Info(qq)
// 		if i == 0 { // 一問目を保存
// 			firstDTO = models.NextQuestion{
// 				QuizID:         qEnt.ID,
// 				QuestionNumber: qq.QuestionNumber,
// 				WordName:       qq.WordName,
// 				ChoicesJpms:    qq.ChoicesJpms,
// 			}
// 		}
// 	}

// 	/* -------------------------------------------------
// 	④ 正常終了レスポンスをセット
// 	-------------------------------------------------*/
// 	resp = &models.CreateQuizResponse{
// 		QuizID:               qEnt.ID,
// 		TotalCreatedQuestion: req.QuestionCount,
// 		NextQuestion:         firstDTO,
// 	}
// 	return // ★ defer が Commit
// }
