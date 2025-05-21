package converter

// import (
// 	"errors"
// 	"word_app/backend/src/models"

// 	"github.com/gin-gonic/gin"
// )

// func BindCreateWord(c *gin.Context) (*models.CreateWordReq, error) {
// 	var req models.CreateWordReq
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		return nil, err
// 	}
// 	uid, ok := c.Get("userID")
// 	if !ok {
// 		return nil, errors.New("unauthorized")
// 	}
// 	req.UserID = uid.(int)
// 	return &req, nil
// }
