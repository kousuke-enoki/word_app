// infrastructure/ratelimit/memory_limiter.go
package ratelimit

import (
	"context"
	"encoding/json"
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

	// 修正: インクリメント前に上限チェック
	// 現在のカウントで判定し、上限に達していればインクリメントしない
	if v.count >= MaxRequests {
		return &RateLimitResult{
			Allowed:      v.last != nil, // last があれば Allowed として返させる仕様に合わせる
			LastPayload:  v.last,
			RetryAfter:   WindowSeconds - int(now.Sub(windowStart).Seconds()),
			CurrentCount: v.count, // インクリメント前のカウント
			WindowStart:  windowStart,
		}, nil
	}

	// 許可された場合のみインクリメント
	v.count++

	return &RateLimitResult{
		Allowed:      true,
		CurrentCount: v.count, // インクリメント後のカウント
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

// ClearCacheForUser は指定されたユーザーIDに関連するキャッシュを削除する
func (m *MemoryLimiter) ClearCacheForUser(ctx context.Context, userID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// すべてのエントリをチェックして、last_resultに指定されたuser_idを含むものを削除
	for _, v := range m.data {
		if v.last == nil {
			continue
		}

		// JSONからuser_idを抽出
		var result struct {
			UserID int `json:"user_id"`
		}
		if err := json.Unmarshal(v.last, &result); err != nil {
			// JSONパースエラーは無視
			continue
		}

		// user_idが一致する場合は、lastだけを削除（countは維持）
		if result.UserID == userID {
			v.last = nil // last_resultだけを削除、countは維持
		}
	}

	return nil
}
