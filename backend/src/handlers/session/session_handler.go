package handler

import (
	"net/http"

	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct{}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{}
}

func SessionStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := contextutil.GetUserRoles(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthenticated",
				"isLogin": false,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Authenticated",
			"userID":  roles.UserID,
			"isLogin": true,
			"isAdmin": roles.IsAdmin,
			"isRoot":  roles.IsRoot,
		})
	}
}
