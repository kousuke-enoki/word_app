// ratelimit/repo.go
package ratelimit

import (
	"context"
	"time"
)

// RateLimiter は汎用的なレート制限インターフェース
type RateLimiter interface {
	CheckRateLimit(ctx context.Context, ip, uaHash, route string) (*RateLimitResult, error)
	SaveLastResult(ctx context.Context, ip, uaHash, route string, payload []byte) error
	// ClearCacheForUser は指定されたユーザーIDに関連するキャッシュを削除する
	// テストユーザー削除時に、削除されたユーザーのトークンを含むキャッシュを削除するために使用
	ClearCacheForUser(ctx context.Context, userID int) error
}

// RateLimitResult はレート制限チェックの結果
type RateLimitResult struct {
	Allowed      bool      // リクエストが許可されているか
	LastPayload  []byte    // 直近の成功レスポンス（JSON）
	RetryAfter   int       // 秒（429時）
	CurrentCount int       // 現在のカウント
	WindowStart  time.Time // ウィンドウ開始時刻
}
