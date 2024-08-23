package src

import (
	"github.com/gin-gonic/gin"
	"eng_app/ent"
	"eng_app/src/handlers"
	"eng_app/src/handlers/user"
)

func SetupRouter(router *gin.Engine, client *ent.Client) {
	router.Use(CORSMiddleware())

	router.GET("/", handlers.RootHandler)
	router.POST("users/sign_up", user.SignUpHandler(client))
	router.POST("users/sign_in", user.SignInHandler(client))
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
