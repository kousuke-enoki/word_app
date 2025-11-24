#!/bin/bash

set -e

# 環境変数で設定可能（デフォルトはlocalhost、コンテナ内ではlocalstack）
LOCALSTACK_ENDPOINT="${LOCALSTACK_ENDPOINT:-http://localhost:4566}"
AWS_REGION="${AWS_REGION:-ap-northeast-1}"
TABLE_NAME="${TABLE_NAME:-rate_limits}"

echo "=========================================="
echo "LocalStack DynamoDB 初期化スクリプト"
echo "=========================================="
echo "Endpoint: $LOCALSTACK_ENDPOINT"
echo "Region: $AWS_REGION"
echo "Table: $TABLE_NAME"
echo "=========================================="

echo ""
echo "⏳ LocalStackの起動を待機中..."
# ヘルスチェックエンドポイントを使用（curlが使える場合）
if command -v curl &> /dev/null; then
  until curl -s "${LOCALSTACK_ENDPOINT}/_localstack/health" | grep -q '"dynamodb": "available"'; do
    echo "   LocalStackがまだ起動していません。5秒後に再試行..."
    sleep 5
  done
else
  # curlがない場合はAWS CLIで直接確認
  until aws dynamodb list-tables --endpoint-url "$LOCALSTACK_ENDPOINT" --region "$AWS_REGION" &>/dev/null; do
    echo "   LocalStackがまだ起動していません。5秒後に再試行..."
    sleep 5
  done
fi
echo "✅ LocalStackが起動しました！"

echo ""
echo "📋 テーブル '$TABLE_NAME' を作成中..."

if aws dynamodb describe-table \
    --table-name "$TABLE_NAME" \
    --endpoint-url "$LOCALSTACK_ENDPOINT" \
    --region "$AWS_REGION" \
    --no-cli-pager \
    2>/dev/null; then
  echo "ℹ️  テーブル '$TABLE_NAME' は既に存在します"
else
  aws dynamodb create-table \
      --table-name "$TABLE_NAME" \
      --attribute-definitions \
          AttributeName=pk,AttributeType=S \
          AttributeName=sk,AttributeType=S \
      --key-schema \
          AttributeName=pk,KeyType=HASH \
          AttributeName=sk,KeyType=RANGE \
      --billing-mode PAY_PER_REQUEST \
      --endpoint-url "$LOCALSTACK_ENDPOINT" \
      --region "$AWS_REGION" \
      --no-cli-pager

  echo "✅ テーブル '$TABLE_NAME' を作成しました"
fi

echo ""
echo "⏱️  TTL設定を適用中..."

aws dynamodb update-time-to-live \
    --table-name "$TABLE_NAME" \
    --time-to-live-specification "Enabled=true, AttributeName=ttl" \
    --endpoint-url "$LOCALSTACK_ENDPOINT" \
    --region "$AWS_REGION" \
    --no-cli-pager \
    2>/dev/null || echo "ℹ️  TTL設定は既に適用されています"

echo ""
echo "=========================================="
echo "✨ LocalStack DynamoDB 初期化完了！"
echo "=========================================="
echo ""
echo "📊 テーブル情報:"
aws dynamodb describe-table \
    --table-name "$TABLE_NAME" \
    --endpoint-url "$LOCALSTACK_ENDPOINT" \
    --region "$AWS_REGION" \
    --no-cli-pager \
    --query 'Table.[TableName,TableStatus,ItemCount,KeySchema]' \
    --output table 2>/dev/null || echo "⚠️  テーブル情報の取得に失敗しました（テーブルは作成済みです）"

echo ""
echo "🚀 準備完了！アプリケーションを起動できます"
