// Package middleware provides HTTP access logging middleware for Gin.
package middleware

import (
	"strings"
	"time"

	"word_app/backend/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const headerReqID = "X-Request-Id"

// AccessLogOpts configures the access log middleware behavior.
type AccessLogOpts struct {
	HealthPath    string // Path to health endpoint (e.g., "/health")
	ExcludeHealth bool   // If true, health endpoint logs are downgraded to debug
}

// AccessLog returns a Gin middleware that logs HTTP access in structured JSON format.
// It handles request ID propagation, severity level determination, and health endpoint suppression.
func AccessLog(logger *logrus.Logger, opts AccessLogOpts) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Request ID決定（X-Request-Idヘッダ優先、なければUUID生成）
		rid := c.GetHeader(headerReqID)
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set("request_id", rid)
		c.Writer.Header().Set(headerReqID, rid)

		// 2. 開始時刻保持
		start := time.Now()

		// 3. 処理実行
		c.Next()

		// 4. 集計
		dur := time.Since(start).Milliseconds()
		status := c.Writer.Status()

		// 5. 重大度判定
		level := logrus.InfoLevel
		switch {
		case status >= 500:
			level = logrus.ErrorLevel
		case status >= 400:
			level = logrus.WarnLevel
		}

		// 6. /health抑制（LOG_EXCLUDE_HEALTH=trueかつinfoレベルの場合、debugにダウングレード）
		path := c.Request.URL.Path
		if opts.ExcludeHealth && opts.HealthPath != "" && strings.HasPrefix(path, opts.HealthPath) && level == logrus.InfoLevel {
			logger.WithFields(baseFields(c, rid, dur)).Debug("access")
			return
		}

		// 7. 出力（1行JSON、msg:"access"固定）
		entry := logger.WithFields(baseFields(c, rid, dur))
		switch level {
		case logrus.ErrorLevel:
			entry.Error("access")
		case logrus.WarnLevel:
			entry.Warn("access")
		default:
			entry.Info("access")
		}
	}
}

// baseFields collects all standard fields for access log entry.
func baseFields(c *gin.Context, rid string, dur int64) logrus.Fields {
	// user_id取得（認証ミドルウェアで設定される）
	var uid any
	if v, ok := c.Get("user_id"); ok {
		uid = v
	} else {
		uid = nil
	}

	// route取得（未マッチ時は"-"）
	route := c.FullPath()
	if route == "" {
		route = "-"
	}

	// リクエスト/レスポンスサイズ
	reqBytes := c.Request.ContentLength
	if reqBytes < 0 {
		reqBytes = 0
	}
	respBytes := c.Writer.Size()

	// エラー分類
	errorKind := errors.Classify(c)

	return logrus.Fields{
		"ts": time.Now().UTC().Format(time.RFC3339Nano),
		// "msg"はentry.Info("access")の引数として使用するため、Fieldsには含めない
		"request_id": rid,
		"route":      route,
		"path":       c.Request.URL.Path,
		"method":     c.Request.Method,
		"status":     c.Writer.Status(),
		"latency_ms": dur,
		"ip":         c.ClientIP(),
		"user_id":    uid,
		"user_agent": c.Request.UserAgent(),
		"referer":    c.Request.Referer(),
		"proto":      c.Request.Proto,
		"host":       c.Request.Host,
		"req_bytes":  reqBytes,
		"resp_bytes": respBytes,
		"error_kind": errorKind,
	}
}
