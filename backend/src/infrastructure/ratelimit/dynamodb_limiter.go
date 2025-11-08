// ratelimit/check_rate_limit.go
package ratelimit

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	RateLimitTableName = "rate_limits"
	WindowSeconds      = 60  // 1分
	MaxRequests        = 5   // 1分以内の上限
	TTLSeconds         = 120 // TTLは2分
)

type DynamoDBLimiter struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBLimiter(tableName string) (*DynamoDBLimiter, error) {
	endpoint := os.Getenv("AWS_ENDPOINT_URL") // 例: http://localstack:4566
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var client *dynamodb.Client
	if endpoint != "" {
		client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			// v1.27+ なら BaseEndpoint が使える
			o.BaseEndpoint = aws.String(endpoint)
			// LocalStack利用時はリージョンも明示しておく
			o.Region = region
		})
	} else {
		client = dynamodb.NewFromConfig(cfg)
	}

	return &DynamoDBLimiter{client: client, tableName: tableName}, nil
}

// --- helpers ---

func attrNum(av types.AttributeValue) (int, bool) {
	n, ok := av.(*types.AttributeValueMemberN)
	if !ok {
		return 0, false
	}
	v, err := strconv.Atoi(n.Value)
	if err != nil {
		return 0, false
	}
	return v, true
}

// getCurrentCountAndWindow は現在のカウントとウィンドウ状態を取得
func (r *DynamoDBLimiter) getCurrentCountAndWindow(
	ctx context.Context,
	key map[string]types.AttributeValue,
	windowStartRFC string,
) (currentCount int, needsReset bool, item map[string]types.AttributeValue, err error) {
	gi, gerr := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      aws.String(r.tableName),
		Key:            key,
		ConsistentRead: aws.Bool(true),
	})

	if gerr != nil || gi.Item == nil {
		return 0, true, nil, nil
	}

	ws, ok := gi.Item["window_start"].(*types.AttributeValueMemberS)
	if !ok || ws.Value != windowStartRFC {
		return 0, true, gi.Item, nil
	}

	// 同じウィンドウ
	count := 0
	if cnt, ok := attrNum(gi.Item["count"]); ok {
		count = cnt
	}
	return count, false, gi.Item, nil
}

// getCachedResponse はキャッシュされたレスポンスを取得
func getCachedResponse(
	item map[string]types.AttributeValue,
	windowStartRFC string,
	currentCount int,
	windowStart time.Time,
) (*RateLimitResult, bool) {
	if item == nil {
		return nil, false
	}

	ws, ok := item["window_start"].(*types.AttributeValueMemberS)
	if !ok || ws.Value != windowStartRFC {
		return nil, false
	}

	lr, ok := item["last_result"].(*types.AttributeValueMemberS)
	if !ok || lr.Value == "" {
		return nil, false
	}

	return &RateLimitResult{
		Allowed:      true,
		LastPayload:  []byte(lr.Value),
		CurrentCount: currentCount,
		WindowStart:  windowStart,
	}, true
}

// incrementOrResetCounter はカウンターをインクリメントまたはリセット
func (r *DynamoDBLimiter) incrementOrResetCounter(
	ctx context.Context,
	key map[string]types.AttributeValue,
	needsReset bool,
	windowStartRFC string,
	ttl int64,
) (*dynamodb.UpdateItemOutput, error) {
	if needsReset {
		return r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:        aws.String(r.tableName),
			Key:              key,
			UpdateExpression: aws.String("SET #c = :one, #ws=:ws, #t=:ttl"),
			ExpressionAttributeNames: map[string]string{
				"#c":  "count",
				"#ws": "window_start",
				"#t":  "ttl",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":one": &types.AttributeValueMemberN{Value: "1"},
				":ws":  &types.AttributeValueMemberS{Value: windowStartRFC},
				":ttl": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", ttl)},
			},
			ReturnValues: types.ReturnValueUpdatedNew,
		})
	}

	return r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:           aws.String(r.tableName),
		Key:                 key,
		UpdateExpression:    aws.String("SET #c = #c + :inc, #t=:ttl"),
		ConditionExpression: aws.String("#ws = :ws"),
		ExpressionAttributeNames: map[string]string{
			"#c":  "count",
			"#ws": "window_start",
			"#t":  "ttl",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
			":ws":  &types.AttributeValueMemberS{Value: windowStartRFC},
			":ttl": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", ttl)},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})
}

// --- main ---

// CheckRateLimit: インクリメント前にカウントをチェックし、上限に達していれば拒否
func (r *DynamoDBLimiter) CheckRateLimit(
	ctx context.Context,
	ip string,
	uaHash string,
	route string,
) (*RateLimitResult, error) {
	now := time.Now().UTC()
	windowStart := now.Truncate(time.Duration(WindowSeconds) * time.Second)
	windowStartRFC := windowStart.Format(time.RFC3339)
	ttl := now.Add(time.Duration(TTLSeconds) * time.Second).Unix()

	pk := fmt.Sprintf("ip#%s", ip)
	sk := fmt.Sprintf("route#%s#ua#%s", route, uaHash)

	key := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"sk": &types.AttributeValueMemberS{Value: sk},
	}

	// 現在のカウントとウィンドウ状態を取得
	currentCount, needsReset, item, err := r.getCurrentCountAndWindow(ctx, key, windowStartRFC)
	if err != nil {
		return nil, fmt.Errorf("failed to get current count: %w", err)
	}

	// 上限チェック
	if currentCount >= MaxRequests {
		// キャッシュされたレスポンスを確認
		if result, found := getCachedResponse(item, windowStartRFC, currentCount, windowStart); found {
			return result, nil
		}

		// キャッシュなし → 429
		remaining := WindowSeconds - int(now.Sub(windowStart).Seconds())
		if remaining <= 0 {
			remaining = WindowSeconds
		}
		return &RateLimitResult{
			Allowed:      false,
			RetryAfter:   remaining,
			CurrentCount: currentCount,
			WindowStart:  windowStart,
		}, nil
	}

	// カウンターをインクリメント
	out, err := r.incrementOrResetCounter(ctx, key, needsReset, windowStartRFC, ttl)
	if err != nil {
		return nil, fmt.Errorf("update rate limit failed: %w", err)
	}

	newCount, ok := attrNum(out.Attributes["count"])
	if !ok {
		newCount = 1
	}

	return &RateLimitResult{
		Allowed:      true,
		CurrentCount: newCount,
		WindowStart:  windowStart,
	}, nil
}

// SaveLastResult: 成功レスポンスを last_result として保存（窓は最新で上書き）
func (r *DynamoDBLimiter) SaveLastResult(
	ctx context.Context,
	ip string,
	uaHash string,
	route string,
	payload []byte,
) error {
	now := time.Now().UTC()
	windowStart := now.Truncate(time.Duration(WindowSeconds) * time.Second).Format(time.RFC3339)
	ttl := now.Add(time.Duration(TTLSeconds) * time.Second).Unix()

	pk := fmt.Sprintf("ip#%s", ip)
	sk := fmt.Sprintf("route#%s#ua#%s", route, uaHash)

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("SET #lr=:lr, #ws=:ws, #t=:ttl"),
		ExpressionAttributeNames: map[string]string{
			"#lr": "last_result",
			"#ws": "window_start",
			"#t":  "ttl",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":lr":  &types.AttributeValueMemberS{Value: string(payload)},
			":ws":  &types.AttributeValueMemberS{Value: windowStart},
			":ttl": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", ttl)},
		},
	})
	return err
}
