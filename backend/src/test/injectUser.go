package test

import "github.com/gin-gonic/gin"

// userID とロールをコンテキストに埋め込む
func InjectUser(c *gin.Context, id int, isRoot bool) {
	c.Set("userID", id)
	c.Set("isAdmin", false)
	c.Set("isRoot", isRoot)
}
