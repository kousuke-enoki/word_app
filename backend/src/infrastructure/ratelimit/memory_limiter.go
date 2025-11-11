// infrastructure/ratelimit/memory_limiter.go
package ratelimit

import (
	"context"
	"sync"
	"time"
)

type key struct{ ip, ua, route, window string }

type MemoryLimiter struct {
	mu   sync.Mutex
	data map[key]*struct {
		count int
		last  []byte // JSON payload
	}
}

func NewMemoryLimiter() *MemoryLimiter {
	return &MemoryLimiter{data: make(map[key]*struct {
		count int
		last  []byte
	})}
}

func (m *MemoryLimiter) CheckRateLimit(ctx context.Context, ip, uaHash, route string) (*RateLimitResult, error) {
	now := time.Now().UTC()
	windowStart := now.Truncate(time.Duration(WindowSeconds) * time.Second)
	k := key{ip: ip, ua: uaHash, route: route, window: windowStart.Format(time.RFC3339)}

	m.mu.Lock()
	defer m.mu.Unlock()

	// ひとつ前の窓は無視（TTL相当の自然消滅。簡易化）
	v, ok := m.data[k]
	if !ok {
		v = &struct {
			count int
			last  []byte
		}{count: 0}
		m.data[k] = v
	}
	v.count++

	// 5回目で上限に達する
	if v.count >= MaxRequests {
		return &RateLimitResult{
			Allowed:      v.last != nil, // last があれば Allowed として返させる仕様に合わせる
			LastPayload:  v.last,
			RetryAfter:   WindowSeconds - int(now.Sub(windowStart).Seconds()),
			CurrentCount: v.count,
			WindowStart:  windowStart,
		}, nil
	}
	return &RateLimitResult{
		Allowed:      true,
		CurrentCount: v.count,
		WindowStart:  windowStart,
		LastPayload:  v.last,
	}, nil
}

func (m *MemoryLimiter) SaveLastResult(ctx context.Context, ip, uaHash, route string, payload []byte) error {
	windowStart := time.Now().UTC().Truncate(time.Duration(WindowSeconds) * time.Second)
	k := key{ip: ip, ua: uaHash, route: route, window: windowStart.Format(time.RFC3339)}

	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[k]
	if !ok {
		v = &struct {
			count int
			last  []byte
		}{count: 0}
		m.data[k] = v
	}
	v.last = payload
	return nil
}
