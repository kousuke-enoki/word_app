package result

import (
	"word_app/backend/src/interfaces"
)

type ResultHandler struct {
	resultService interfaces.ResultService
}

func NewResultHandler(resultService interfaces.ResultService) *ResultHandler {
	return &ResultHandler{resultService: resultService}
}
