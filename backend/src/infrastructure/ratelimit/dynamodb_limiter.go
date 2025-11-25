// ratelimit/check_rate_limit.go
package ratelimit

import (
	"context"
	"encoding/json"
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
)

var (
	// 環境変数から読み込む（デフォルト値付き）
	WindowSeconds = getEnvInt("RATE_LIMIT_WINDOW_SEC", 60)  // 1分
	MaxRequests   = getEnvInt("RATE_LIMIT_MAX_REQUESTS", 1) // 1分以内の上限（既定: 1回）
	TTLSeconds    = getEnvInt("RATE_LIMIT_TTL_SEC", 120)    // TTLは2分（ウィンドウ + 余裕）
)

// getEnvInt は環境変数を整数として取得（デフォルト値付き）
func getEnvInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}

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

// buildRateLimitKey は要件仕様に従ったキーを構築
// PK: "RL#" + route + "#" + ip + "#" + uaHash
// SK: window_start (ISO-8601形式)
func buildRateLimitKey(ip, uaHash, route, windowStartRFC string) map[string]types.AttributeValue {
	pk := fmt.Sprintf("RL#%s#%s#%s", route, ip, uaHash)
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"sk": &types.AttributeValueMemberS{Value: windowStartRFC},
	}
}

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

	// DynamoDBのエラーは即座に返す（fail-closed）
	if gerr != nil {
		return 0, false, nil, fmt.Errorf("failed to get rate limit item: %w", gerr)
	}

	// レコードが存在しない場合は正常なケース（初回リクエストなど）
	if gi.Item == nil {
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

	// TTL計算: window_start + window + 余裕バッファ(5秒)
	ttl := windowStart.Add(time.Duration(WindowSeconds)*time.Second + 5*time.Second).Unix()

	// 要件仕様のキー設計を使用: PK="RL#route#ip#uaHash", SK=window_start
	key := buildRateLimitKey(ip, uaHash, route, windowStartRFC)

	// 現在のカウントとウィンドウ状態を取得
	currentCount, needsReset, item, err := r.getCurrentCountAndWindow(ctx, key, windowStartRFC)
	if err != nil {
		return nil, fmt.Errorf("failed to get current count: %w", err)
	}

	// 上限チェック
	if currentCount >= MaxRequests {
		// RetryAfter時間を計算
		remaining := WindowSeconds - int(now.Sub(windowStart).Seconds())
		if remaining <= 0 {
			remaining = WindowSeconds
		}

		// キャッシュされたレスポンスを確認
		if result, found := getCachedResponse(item, windowStartRFC, currentCount, windowStart); found {
			// レート制限超過時なので、Allowedをfalseにし、RetryAfterを設定
			result.Allowed = false
			result.RetryAfter = remaining
			return result, nil
		}

		// キャッシュなし → 429
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
	windowStart := now.Truncate(time.Duration(WindowSeconds) * time.Second)
	windowStartRFC := windowStart.Format(time.RFC3339)

	// TTL計算: window_start + window + 余裕バッファ(5秒)
	ttl := windowStart.Add(time.Duration(WindowSeconds)*time.Second + 5*time.Second).Unix()

	// 要件仕様のキー設計を使用: PK="RL#route#ip#uaHash", SK=window_start
	key := buildRateLimitKey(ip, uaHash, route, windowStartRFC)

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(r.tableName),
		Key:              key,
		UpdateExpression: aws.String("SET #lr=:lr, #ws=:ws, #t=:ttl"),
		ExpressionAttributeNames: map[string]string{
			"#lr": "last_result",
			"#ws": "window_start",
			"#t":  "ttl",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":lr":  &types.AttributeValueMemberS{Value: string(payload)},
			":ws":  &types.AttributeValueMemberS{Value: windowStartRFC},
			":ttl": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", ttl)},
		},
	})
	return err
}

// ClearCacheForUser は指定されたユーザーIDに関連するキャッシュを削除する
// last_resultのJSONからuser_idを抽出し、一致するキャッシュを削除
func (r *DynamoDBLimiter) ClearCacheForUser(ctx context.Context, userID int) error {
	// DynamoDBのすべてのアイテムをスキャンし、last_resultに指定されたuser_idを含むものを削除
	// 注意: 全スキャンは効率が悪いが、テストユーザー削除は頻繁でないため許容
	// 将来的にはGSIでuser_idで検索できるようにすると効率的

	// Scanで全アイテムを取得
	paginator := dynamodb.NewScanPaginator(r.client, &dynamodb.ScanInput{
		TableName:      aws.String(r.tableName),
		ConsistentRead: aws.Bool(true),
	})

	var itemsToUpdate []map[string]types.AttributeValue

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to scan cache items: %w", err)
		}

		// 各アイテムのlast_resultをチェック
		for _, item := range output.Items {
			lr, ok := item["last_result"].(*types.AttributeValueMemberS)
			if !ok || lr.Value == "" {
				continue
			}

			// JSONからuser_idを抽出
			var result struct {
				UserID int `json:"user_id"`
			}
			if err := json.Unmarshal([]byte(lr.Value), &result); err != nil {
				// JSONパースエラーは無視（古い形式の可能性）
				continue
			}

			// user_idが一致する場合は更新対象に追加
			if result.UserID == userID {
				// キーを構築（PKとSKのみ）
				key := map[string]types.AttributeValue{
					"pk": item["pk"],
					"sk": item["sk"],
				}
				itemsToUpdate = append(itemsToUpdate, key)
			}
		}
	}

	// last_resultだけを削除（countは維持する）
	for _, key := range itemsToUpdate {
		_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:        aws.String(r.tableName),
			Key:              key,
			UpdateExpression: aws.String("REMOVE #lr"), // last_resultだけを削除
			ExpressionAttributeNames: map[string]string{
				"#lr": "last_result",
			},
		})
		if err != nil {
			// ログは出すが、一部失敗しても続行
			// 本番環境では構造化ログを使用
			continue
		}
	}

	return nil
}
