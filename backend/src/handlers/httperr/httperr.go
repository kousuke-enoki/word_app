// interfaces/http/httperr/httperr.go
package httperr

import (
	"errors"
	"net/http"

	"word_app/backend/logger/logx"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Kind→HTTP の一元マッピング＋薄いラッパー
func Write(c *gin.Context, err error) {
	ctx := c.Request.Context()
	var ae *apperror.Error
	if !errors.As(err, &ae) {
		logx.From(ctx).WithError(err).Error("unhandled error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// 構造化ログ（メタと原因は全部ログへ）
	logger := logx.From(ctx).WithFields(logrus.Fields{
		"kind":  ae.Kind,
		"meta":  ae.Meta,
		"error": ae.Message,
	})
	if ae.Err != nil {
		logger = logger.WithError(ae.Err)
	}

	// レベルの方針（例）
	switch ae.Kind {
	case apperror.Validation, apperror.InvalidCredential, apperror.Conflict, apperror.NotFound:
		logger.Warn("request error")
	default:
		logger.Error("server error")
	}

	status := StatusOf(ae.Kind)
	body := gin.H{"error": ae.Message}
	if f, ok := ae.Meta["fields"]; ok {
		body["fields"] = f
	}
	c.JSON(status, body)
}

func StatusOf(k apperror.Kind) int {
	switch k {
	case apperror.Unauthorized:
		return http.StatusUnauthorized
	case apperror.Forbidden:
		return http.StatusForbidden
	case apperror.NotFound:
		return http.StatusNotFound
	case apperror.Conflict:
		return http.StatusConflict
	case apperror.Validation, apperror.InvalidCredential:
		return http.StatusBadRequest
	case apperror.TooManyRequests:
		return http.StatusTooManyRequests
	case apperror.TooLargeRequests:
		return http.StatusRequestEntityTooLarge
	case apperror.BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
