package models

// type CreateQuizRequest struct {
// 	QuizCount       int    `json:"questionCount"`
// 	TargetWords     string `json:"targetWordTypes"`
// 	PartsOfSpeeches []int  `json:"partsOfSpeeches"`
// 	UserID          int    `json:"userId"`
// }

type CreateQuizReq struct {
	QuestionCount       int   `json:"questionCount" binding:"required,min=10,max=100"`
	IsSaveResult        bool  `json:"isSaveResult"`
	IsRegisteredWords   int   `json:"isRegisteredWords" binding:"oneof=0 1 2"`
	CorrectRate         int   `json:"correctRate" binding:"required,min=0,max=100"`
	AttentionLevelList  []int `json:"attentionLevelList"` // 1‑5 配列
	PartsOfSpeeches     []int `json:"partsOfSpeeches"`    // 1‑12 配列
	IsIdioms            int   `json:"isIdioms" binding:"oneof=0 1 2"`
	IsSpecialCharacters int   `json:"isSpecialCharacters" binding:"oneof=0 1 2"`
}

type CreateQuizResponse struct {
	QuizID               int          `json:"quizID"`
	TotalCreatedQuestion int          `json:"totalCreatedQuestion"`
	NextQuestion         NextQuestion `json:"nextQuestion"`
}

type PostAnswerQuestionRequest struct {
	QuizID         int `json:"quizID"`
	QuestionNumber int `json:"questionNumber"`
	AnswerJpmID    int `json:"answerJpmId"`
}

// type PostAnswerQuestionResponse struct {
// 	IsCorrectBefore bool `json:"isCorrectBefore"`
// }

type GetQuizRequest struct {
	QuizID               *int `json:"quizId,omitempty" form:"quizId"`
	BeforeQuestionNumber *int `json:"questionNumber,omitempty" form:"questionNumber"`
}

type GetQuizResponse struct {
	IsRunningQuiz bool         `json:"isRunningQuiz"`
	NextQuestion  NextQuestion `json:"nextQuestion"`
}

type AnswerRouteRes struct {
	IsFinish  bool `json:"isFinish"`
	IsCorrect bool `json:"isCorrect"`
	// PostAnswerQuestionResponse PostAnswerQuestionResponse `json:"postAnswerQuestionResponse"`
	NextQuestion NextQuestion `json:"nextQuestion"`
	Result       Result       `json:"result"`
}

type NextQuestion struct {
	QuizID         int         `json:"quizId"`
	QuestionNumber int         `json:"questionNumber"`
	WordName       string      `json:"wordName"`
	ChoicesJpms    []ChoiceJpm `json:"choicesJpms"`
}

type ChoiceJpm struct {
	JapaneseMeanID int    `json:"japaneseMeanID"`
	Name           string `json:"name"`
}

type Result struct {
	QuizNumber          int              `json:"quizNumber"`
	TotalQuestionsCount int              `json:"totalQuestionsCount"`
	CorrectCount        int              `json:"correctCount"`
	ResultCorrectRate   float64          `json:"resultCorrectRate"`
	ResultSetting       ResultSetting    `json:"ResultSetting"`
	ResultQuestions     []ResultQuestion `json:"resultQuestions"`
}

type ResultQuestion struct {
	QuestionNumber int         `json:"questionNumber"`
	WordName       string      `json:"wordName"`
	PosID          int         `json:"posID"`
	CorrectJpmId   int         `json:"correctJpmId"`
	ChoicesJpms    []ChoiceJpm `json:"choicesJpms"`
	AnswerJpmId    int         `json:"answerJpmId"`
	IsCorrect      bool        `json:"isCorrect"`
	TimeMs         int         `json:"timeMs"`
}

type ResultSetting struct {
	IsSaveResult        bool  `json:"isSaveResult"`
	IsRegisteredWords   int   `json:"isRegisteredWords"`
	SettingCorrectRate  int   `json:"settingCorrectRate"`
	IsIdioms            int   `json:"isIdioms"`
	IsSpecialCharacters int   `json:"isSpecialCharacters"`
	AttentionLevelList  []int `json:"attentionLevelList"`
	ChoicesPosIds       []int `json:"choicesPosIds"`
}
