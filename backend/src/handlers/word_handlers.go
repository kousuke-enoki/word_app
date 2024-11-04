package handlers

import (
	"word_app/backend/ent"

	"github.com/gin-gonic/gin"
)

type wordHandler struct {
	client *ent.Client
}

func NewWordHandler(client *ent.Client) *wordHandler {
	return &wordHandler{client: client}
}

func (h *wordHandler) AllWordList(c *gin.Context) {
	// 単語リストの取得処理
}

func (h *wordHandler) WordShow(c *gin.Context) {
	// 単語の詳細表示処理
}
