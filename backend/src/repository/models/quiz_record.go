package models

type QuizRecord struct {
	QuestionCount      int
	IsSaveResult       bool
	IsRegisteredWords  int
	SettingCorrectRate int
	IsIdioms           int
	IsSpecialChars     int
	AttentionLevels    []int
	ChoicePosIDs       []int
}

type QuizQuestionRecord struct {
	QuestionNumber int
	WordID         int
	PosID          int
	CorrectJpmID   int
	ChoicesJSON    []byte
	WordName       string // NextQuestionDTO 用に渡したい場合だけ
}
