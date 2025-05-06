package quiz_service

// func (s *QuizServiceImpl) FinishQuizAsdf(
// 	ctx context.Context,
// 	userID int,
// 	req *models.GetQuizRequest,
// ) (*models.ResultRes, error) {

// 	tx, err := s.client.Tx(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if p := recover(); p != nil {
// 			_ = tx.Rollback()
// 			panic(p)
// 		} else if err != nil {
// 			_ = tx.Rollback()
// 		} else {
// 			err = tx.Commit()
// 		}
// 	}()

// 	// ─── 1. クイズ本体 + 質問を取得 ─────────────────────────────
// 	q, err := tx.Quiz.
// 		Query().
// 		Where(
// 			quiz.IDEQ(req.QuizID),
// 			quiz.UserID(userID),
// 		).
// 		WithQuizQuestions().
// 		Only(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// ─── 2. 正誤集計 & DTO 構築 ─────────────────────────────
// 	var (
// 		correctCnt      int
// 		resultQuestions = make([]models.ResultQuestion, 0, len(q.Edges.QuizQuestions))
// 	)

// 	for _, qq := range q.Edges.QuizQuestions {
// 		// nil セーフガード
// 		isCor := qq.IsCorrect != nil && *qq.IsCorrect
// 		if isCor {
// 			correctCnt++
// 		}

// 		resQ := models.ResultQuestion{
// 			QuestionNumber: qq.QuestionNumber,
// 			WordName:       qq.WordName,
// 			PosID:          qq.PosID,
// 			CorrectJpmId:   qq.CorrectJpmID,
// 			ChoicesJpms:    qq.ChoicesJpms,
// 			AnswerJpmId:    derefInt(qq.AnswerJpmID),
// 			IsCorrect:      isCor,
// 			TimeMs:         derefInt(qq.TimeMs),
// 		}
// 		resultQuestions = append(resultQuestions, resQ)
// 	}

// 	// 小数点込みで正答率を計算
// 	cRate := float64(correctCnt) / float64(q.TotalQuestionsCount)

// 	// ─── 3. クイズを更新 ────────────────────────────────────
// 	_, err = tx.Quiz.
// 		UpdateOne(q).
// 		SetIsRunning(false).
// 		SetCorrectCount(correctCnt).
// 		SetResultCorrectRate(cRate).
// 		Save(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// ─── 4. isSaveResult == false → soft delete ───────────
// 	if !q.IsSaveResult {
// 		now := time.Now()
// 		if _, err = tx.Quiz.
// 			UpdateOne(q).
// 			SetDeletedAt(now).
// 			Save(ctx); err != nil {
// 			return nil, err
// 		}

// 		if _, err = tx.QuizQuestion.
// 			Update().
// 			Where(quizquestion.QuizIDEQ(q.ID)).
// 			SetDeletedAt(now).
// 			Save(ctx); err != nil {
// 			return nil, err
// 		}
// 	}

// 	// ─── 5. レスポンス組立て ───────────────────────────────
// 	return &models.ResultRes{
// 		QuizNumber:          q.QuizNumber,
// 		TotalQuestionsCount: q.TotalQuestionsCount,
// 		CorrectCount:        correctCnt,
// 		ResultCorrectRate:   cRate,
// 		ResultSetting: models.ResultSetting{
// 			IsSaveResult:        q.IsSaveResult,
// 			IsRegisteredWords:   q.IsRegisteredWords,
// 			SettingCorrectRate:  q.SettingCorrectRate,
// 			IsIdioms:            q.IsIdioms,
// 			IsSpecialCharacters: q.IsSpecialCharacters,
// 			AttentionLevelList:  q.AttentionLevelList,
// 			ChoicesPosIds:       q.ChoicesPosIds,
// 		},
// 		ResultQuestions: resultQuestions,
// 	}, nil
// }

// // ヘルパ
// func derefInt(ptr *int) int {
// 	if ptr == nil {
// 		return 0
// 	}
// 	return *ptr
// }
