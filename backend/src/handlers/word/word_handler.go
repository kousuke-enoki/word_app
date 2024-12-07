package word

import (
	"word_app/backend/src/interfaces"
)

type WordHandler struct {
	wordService interfaces.WordService
}

func NewWordHandler(wordService interfaces.WordService) *WordHandler {
	return &WordHandler{wordService: wordService}
}
