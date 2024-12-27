package exam

import (
	"word_app/backend/src/interfaces"
)

type ExamHandler struct {
	examService interfaces.ExamService
}

func NewExamHandler(examService interfaces.ExamService) *ExamHandler {
	return &ExamHandler{examService: examService}
}
