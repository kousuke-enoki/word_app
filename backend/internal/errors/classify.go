// Package errors provides error classification utilities for access logging.
package errors

import (
	"context"
	"errors"
	"net"

	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Classify analyzes errors in gin.Context and returns an error kind string.
// Priority order: timeout > validation > db > external > unknown
// Returns empty string if no errors found.
func Classify(c *gin.Context) string {
	if len(c.Errors) == 0 {
		return ""
	}

	// 複数エラー時は最重要1件を返す（優先順位順にチェック）
	for _, e := range c.Errors {
		if kind := classifyError(e.Err); kind != "" {
			return kind
		}
	}

	// その他のエラー
	return "unknown"
}

// classifyError classifies a single error and returns its kind.
// Returns empty string if the error doesn't match any known category.
func classifyError(err error) string {
	// タイムアウトエラー（直接チェック）
	if isTimeoutError(err) {
		return "timeout"
	}

	// apperror.Errorにラップされている場合
	var ae *apperror.Error
	if errors.As(err, &ae) {
		// バリデーションエラー
		if ae.Kind == apperror.Validation {
			return "validation"
		}

		// ラップされたエラーを再帰的にチェック
		if ae.Err != nil {
			if isTimeoutError(ae.Err) {
				return "timeout"
			}
			// DBエラー（PostgreSQL）
			var pqErr *pq.Error
			if errors.As(ae.Err, &pqErr) {
				return "db"
			}
		}
	}

	// DBエラー（PostgreSQL）- apperror.Errorにラップされていない場合
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return "db"
	}

	return ""
}

// isTimeoutError checks if an error is a timeout error.
func isTimeoutError(err error) bool {
	if err == context.DeadlineExceeded || err == context.Canceled {
		return true
	}
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		return true
	}
	return false
}
