package word

import (
	"word_app/backend/src/interfaces/http/word"
)

type Handler struct {
	wordService word.Service
}

func NewHandler(
	wordService word.Service,
) *Handler {
	return &Handler{
		wordService: wordService,
	}
}
