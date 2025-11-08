// infrastructure/ratelimit/factory.go
package ratelimit

import (
	"fmt"
	"os"
)

func NewRateLimiterFromEnv() (RateLimiter, error) {
	switch os.Getenv("RATE_LIMIT_BACKEND") {
	case "dynamodb":
		table := os.Getenv("RATE_LIMIT_TABLE")
		if table == "" {
			return nil, fmt.Errorf("RATE_LIMIT_TABLE is required")
		}
		repo, err := NewDynamoDBLimiter(table) // 既存相当
		if err != nil {
			return nil, err
		}
		return repo, nil
	case "memory", "":
		return NewMemoryLimiter(), nil
	default:
		return NewNoopLimiter(), nil
	}
}
