// handlers/http_error.go
package handlers

import (
	"net/http"

	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, err error) {
	if ae, ok := err.(*apperror.Error); ok {
		status := mapToStatus(ae.Kind)
		c.JSON(status, gin.H{"error": gin.H{"code": ae.Kind, "message": ae.Message}})
		return
	}
	// 予期しないエラー
	c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": apperror.Internal, "message": "internal error"}})
}

func mapToStatus(k apperror.Kind) int {
	switch k {
	case apperror.Unauthorized:
		return http.StatusUnauthorized
	case apperror.Forbidden:
		return http.StatusForbidden
	case apperror.NotFound:
		return http.StatusNotFound
	case apperror.Conflict:
		return http.StatusConflict
	case apperror.Validation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
