package exam_service

import (
	"errors"
	"word_app/backend/ent"
)

type ExamServiceImpl struct {
	client *ent.Client
}

func NewExamService(client *ent.Client) *ExamServiceImpl {
	return &ExamServiceImpl{client: client}
}

var (
	ErrExamNotFound = errors.New("word not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrDeleteExam   = errors.New("failed to delete word")
)
