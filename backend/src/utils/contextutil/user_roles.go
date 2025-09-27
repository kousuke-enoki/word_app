package contextutil

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type UserRoles struct {
	UserID  int
	IsAdmin bool
	IsRoot  bool
	IsTest  bool
}

// errorを返すのは、userIDがcontextに存在しない場合のみ
// つまり未ログイン状態で呼び出された場合なので、基本的にはログイン必須のAPIでしか使わない
// その場合はtopに戻す
func GetUserRoles(c *gin.Context) (*UserRoles, error) {
	id, ok := c.Get("userID")
	if !ok {
		return nil, errors.New("userID not found in context")
	}
	isAdmin, _ := c.Get("isAdmin")
	isRoot, _ := c.Get("isRoot")
	isTest, _ := c.Get("isTest")

	return &UserRoles{
		UserID:  id.(int),
		IsAdmin: isAdmin.(bool),
		IsRoot:  isRoot.(bool),
		IsTest:  isTest.(bool),
	}, nil
}
