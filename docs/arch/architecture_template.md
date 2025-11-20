# Architecture（貼り付け用テンプレ）

## 構成要素

- **Client**: React（Vercel）
- **API**: API Gateway → AWS Lambda（Go/Gin）
- **DB**: PostgreSQL（Aurora Serverless v2 / Supabase いずれか）
- **Storage**: S3（辞書 JSON 配置—任意）
- **Observability**: 構造化ログ（zap/zerolog）、Request-ID、k6（P0）、（将来：OTel + Grafana）
- **CI/CD**: GitHub Actions（lint/test/build/k6、デプロイ）

## クリーンアーキテクチャ

本アプリケーションは**クリーンアーキテクチャ**を採用。レイヤー構成は以下の通り：

- **Handlers 層**（プレゼンテーション）→ **UseCase 層**（アプリケーション）→ **Domain 層**（エンティティ）← **Infrastructure 層**（実装）
- **移行状況**: Auth/Bulk/Setting/User/JWT は UseCase 層へ移行完了 ✅、Word/Quiz/Result は移行進行中 🔄
- **依存性逆転**: Repository インターフェースは Domain 層、実装は Infrastructure 層

## データフロー（要点）

1. React → `APIGW`（HTTPS, JWT or セッション）
2. `Lambda`（Gin）で認証・ルーティング → `UseCase層`でビジネスロジック実行
3. `UseCase層` → `Repositoryインターフェース`（Domain 層定義）→ `PostgreSQL` にクエリ
4. `words` 検索は prefix インデックスを利用（全文検索は未導入）
5. クイズ生成は登録語から抽出、結果は `quiz_results` に永続化
6. すべての HTTP ハンドラで `request_id` を発行・伝搬（ログ相関）

## 信頼性/運用

- **SLO**: p95 < 200ms, error < 1%（login / search / quiz）
- **Rate Limit**: 未ログイン=IP 10 r/s、認証=User 20 r/s（429/Retry-After）
- **Timeout/Retry**: DB=3s、外部 API=2–5s、指数バックオフ、永続的 4xx は即 fail
- **エラー分類**: `validation|db|external|unknown` を JSON ログへ出力
- **Demo 保護**: 書込上限・CRON で日次リセット（00:00 JST）

## セキュリティ

- JWT（短寿命 Access + 長寿命 Refresh/将来）、`RequireActiveUser`（削除/無効化ユーザは拒否）
- CORS: デモフロントのオリジンのみ許可
- Secrets: 環境変数（`.env.*` は `.env.example` で項目のみ公開）

## パフォーマンス設計

- **インデックス**: `words(name text_pattern_ops)` 等（代表クエリ最適化）
- **N+1 排除**: Ent の `WithXXX` で解消
- **スケール**: Lambda（Provisioned Concurrency=1 でコールドスタート緩和）、DB は自動スケール/接続数監視

## 既知の制約

- Full-text search 未対応（prefix 最適化のみ）
- CDN/キャッシュは最小（API 直）
- マルチリージョン冗長化なし
- Word/Quiz/Result 機能は Service 層経由（UseCase 層への移行予定）

## ポート/プロトコル

- Client→API: HTTPS/443（REST, JSON）
- API→DB: Postgres/5432（VPC or マネージド接続）

## 設計原則

- **Request-ID で全レイヤ追跡**、JSON ログで `route/status/latency_ms/error_kind` を出力
- **Rate Limit と Timeout/Retry** を全 I/O に適用（k6 で SLO を継続検証）
- **Index/N+1 対策**で `GET /words` の p95 を改善（Before/After は Performance 章）
- **依存性逆転**: Repository インターフェースは Domain 層、実装は Infrastructure 層
- **テスタビリティ**: モック（mockery）で各レイヤーを独立テスト可能
