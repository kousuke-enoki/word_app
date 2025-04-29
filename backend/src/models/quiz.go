package models

type CreateQuizRequest struct {
	QuizCount       int    `json:"questionCount"`
	TargetWords     string `json:"targetWordTypes"`
	PartsOfSpeeches []int  `json:"partsOfSpeeches"`
	UserID          int    `json:"userId"`
}

type CreateQuizResponse struct {
	QuizID        int            `json:"quizID"`
	TotalQuizs    string         `json:"totalQuizs"`
	QuizQuestions []QuizQuestion `json:"quizQuestions"`
}

type GetQuizRequest struct {
	QuizID int `json:"quizID"`
	UserID int `json:"userID"`
}

type GetQuizResponse struct {
	QuizID        int            `json:"quizID"`
	TotalQuizs    string         `json:"totalQuizs"`
	QuizQuestions []QuizQuestion `json:"quizQuestions"`
}

type QuizQuestion struct {
	QuizQuestionID int           `json:"quizQuestionID"`
	WordName       string        `json:"wordName"`
	QuestionJpms   []QuestionJpm `json:"questionJpms"`
}

type QuestionJpm struct {
	JapaneseMeanID int    `json:"japaneseMeanID"`
	Name           string `json:"name"`
}
