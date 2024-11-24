package word

import (
	"word_app/backend/src/interfaces"

	"github.com/sirupsen/logrus"
)

type WordHandler struct {
	wordService interfaces.WordService
}

func NewWordHandler(wordService interfaces.WordService) *WordHandler {
	logrus.Info("uuuy")
	return &WordHandler{wordService: wordService}
}
