package word_service

import (
	"word_app/backend/ent"
)

type WordServiceImpl struct {
	client *ent.Client
}

func NewWordService(client *ent.Client) *WordServiceImpl {
	return &WordServiceImpl{client: client}
}
