// interfaces/http/middleware/request_logger.go
package logger

import (
	"time"

	"word_app/backend/logger/logx"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.NewString()
		user, err := contextutil.GetUserRoles(c)
		actorID := 0
		if err == nil && user != nil {
			actorID = user.UserID
		}

		entry := logrus.WithFields(logrus.Fields{
			"request_id": reqID,
			"actor_id":   actorID,
			"path":       c.FullPath(),
			"method":     c.Request.Method,
		})
		ctx := logx.With(c.Request.Context(), entry)
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()

		// アクセスログ
		level := logrus.InfoLevel
		if status >= 500 {
			level = logrus.ErrorLevel
		}
		if status >= 400 && status < 500 {
			level = logrus.WarnLevel
		}

		logx.From(ctx).WithFields(logrus.Fields{
			"status":  status,
			"latency": latency.String(),
		}).Log(level, "access")
	}
}
