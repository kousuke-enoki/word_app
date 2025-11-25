// ratelimit/repo_noop.go
package ratelimit

import (
	"context"
	"time"
)

type NoopLimiter struct{}

func NewNoopLimiter() *NoopLimiter { return &NoopLimiter{} }

func (n *NoopLimiter) CheckRateLimit(ctx context.Context, ip, uaHash, route string) (*RateLimitResult, error) {
	return &RateLimitResult{Allowed: true, CurrentCount: 1, WindowStart: time.Now()}, nil
}

func (n *NoopLimiter) SaveLastResult(ctx context.Context, ip, uaHash, route string, payload []byte) error {
	return nil
}

func (n *NoopLimiter) ClearCacheForUser(ctx context.Context, userID int) error {
	// NoOp実装なので何もしない
	return nil
}
