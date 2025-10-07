// pkg/logx/logx.go
package logx

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ctxKey struct{}

func With(ctx context.Context, e *logrus.Entry) context.Context {
	return context.WithValue(ctx, ctxKey{}, e)
}

func From(ctx context.Context) *logrus.Entry {
	if v := ctx.Value(ctxKey{}); v != nil {
		if e, ok := v.(*logrus.Entry); ok {
			return e
		}
	}
	return logrus.WithContext(ctx)
}
