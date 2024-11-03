// handlers/word_handler.go
package handlers

import (
	"word_app/backend/ent"
	"word_app/backend/src/handlers/word"

	"github.com/gin-gonic/gin"
)

type WordHandler struct {
	client *ent.Client
}

func NewWordHandler(client *ent.Client) *WordHandler {
	return &WordHandler{client: client}
}

func (h *WordHandler) AllWordListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		word.AllWordListHandler(c, h.client)
	}
}

func (h *WordHandler) WordShowHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		word.WordShowHandler(c, h.client)
	}
}
