package models

type CreateExamRequest struct {
	ExamCount       int    `json:"questionCount"`
	TargetWords     string `json:"targetWordTypes"`
	PartsOfSpeeches []int  `json:"partsOfSpeeches"`
	UserID          int    `json:"userId"`
}

type CreateExamResponse struct {
	ExamID        int            `json:"examID"`
	TotalExams    string         `json:"totalExams"`
	ExamQuestions []ExamQuestion `json:"examQuestions"`
}

type GetExamRequest struct {
	ExamID int `json:"examID`
	UserID int `json:"userId"`
}

type GetExamResponse struct {
	ExamID        int            `json:"examID"`
	TotalExams    string         `json:"totalExams"`
	ExamQuestions []ExamQuestion `json:"examQuestions"`
}

type ExamQuestion struct {
	ExamQuestionID int           `json:"examQuestionID"`
	WordName       string        `json:"wordName"`
	QuestionJpms   []QuestionJpm `json:"questionJpms"`
}

type QuestionJpm struct {
	JapaneseMeanID int    `json:"japaneseMeanID"`
	Name           string `json:"name"`
}
