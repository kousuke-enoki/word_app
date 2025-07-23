package result

import "word_app/backend/src/interfaces/http/result"

type Handler struct {
	resultService result.Service
}

func NewHandler(resultService result.Service) *Handler {
	return &Handler{resultService: resultService}
}
