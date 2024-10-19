package handlers

import (
	"github.com/gin-gonic/gin"
)

func RootHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":      "Redirect to root",
		"redirect_url": "/",
	})
}
