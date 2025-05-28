// interface/contextutil/userroles.go
package contextutil

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type UserRoles struct {
	UserID  int
	IsAdmin bool
	IsRoot  bool
}

func GetUserRoles(c *gin.Context) (*UserRoles, error) {
	id, ok := c.Get("userID")
	if !ok {
		return nil, errors.New("userID not found in context")
	}
	isAdmin, _ := c.Get("isAdmin")
	isRoot, _ := c.Get("isRoot")

	return &UserRoles{
		UserID:  id.(int),
		IsAdmin: isAdmin.(bool),
		IsRoot:  isRoot.(bool),
	}, nil
}
