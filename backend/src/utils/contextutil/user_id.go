// src/interface/contextutil/userid.go
package contextutil

import (
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
)

// MustUserID returns the authenticated user ID stored in Gin context.
// It fails with error when not present or of wrong type.
func MustUserID(c *gin.Context) (int, error) {
	v, ok := c.Get("userID")
	if !ok {
		return 0, apperror.Unauthorizedf("unauthorized: userID not found in context", nil)
	}
	id, ok := v.(int)
	if !ok {
		return 0, apperror.Unauthorizedf("unauthorized: userID not found in context", nil)
	}
	return id, nil
}
