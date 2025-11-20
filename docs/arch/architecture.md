# Architecture

## 概要

本アプリケーションは**クリーンアーキテクチャ**を採用し、レイヤー間の依存関係を明確に分離しています。認証・設定・ユーザー管理などのコア機能は既に UseCase 層へ移行済みで、Word/Quiz/Result 機能は移行進行中です。

## 構成要素

- **Client**: React（Vercel）
- **API**: API Gateway → AWS Lambda（Go/Gin）
- **DB**: PostgreSQL（Aurora Serverless v2 / Supabase いずれか）
- **Storage**: S3（辞書 JSON 配置—任意）
- **Observability**: 構造化ログ（zap/zerolog）、Request-ID、k6（P0）、（将来：OTel + Grafana）
- **CI/CD**: GitHub Actions（lint/test/build/k6、デプロイ）

## クリーンアーキテクチャのレイヤー構成

### レイヤー構造

```
┌─────────────────────────────────────┐
│  Handlers (HTTP)                    │  ← プレゼンテーション層
│  - src/handlers/                    │
│  - src/interfaces/http/              │
├─────────────────────────────────────┤
│  UseCase (Business Logic)           │  ← アプリケーション層
│  - src/usecase/                     │
│    - auth, bulk, setting, user, jwt │
├─────────────────────────────────────┤
│  Domain (Entities)                  │  ← ドメイン層
│  - src/domain/                      │
│    - user, root_config, etc.        │
├─────────────────────────────────────┤
│  Infrastructure (External)          │  ← インフラストラクチャ層
│  - src/infrastructure/              │
│    - repository/, jwt/, ratelimit/   │
└─────────────────────────────────────┘
```

### 依存関係の方向

- **Handlers** → **UseCase** → **Domain** ← **Infrastructure**
- 外側のレイヤーは内側のレイヤーに依存（依存性逆転の原則）
- Repository インターフェースは Domain 層に定義、実装は Infrastructure 層

### 移行状況

**✅ 移行完了（UseCase 層）**

- `Auth`: 認証・ログイン・ログアウト
- `Bulk`: 一括登録・トークン化
- `Setting`: 設定取得・更新
- `User`: ユーザー管理
- `JWT`: トークン検証

**🔄 移行進行中（Service 層 → UseCase 層）**

- `Word`: 単語 CRUD・登録語管理（Service 層を経由、UseCase 層への移行予定）
- `Quiz`: クイズ生成（Service 層を経由、UseCase 層への移行予定）
- `Result`: クイズ結果保存（Service 層を経由、UseCase 層への移行予定）

**移行方針**: 既存の Service 層は「薄い Facade」として存続させつつ、段階的に UseCase 層へ移行。DI コンテナ（`internal/di/`）で依存関係を管理。

## データフロー（要点）

1. **React** → `APIGW`（HTTPS, JWT or セッション）
2. **Lambda**（Gin）で認証・ルーティング → **UseCase 層**でビジネスロジック実行
3. **UseCase 層** → **Repository インターフェース**（Domain 層定義）→ **PostgreSQL** にクエリ
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

- **Request-ID で全レイヤ追跡**: JSON ログで `route/status/latency_ms/error_kind` を出力
- **Rate Limit と Timeout/Retry**: 全 I/O に適用（k6 で SLO を継続検証）
- **Index/N+1 対策**: `GET /words` の p95 を改善（Before/After は Performance 章）
- **依存性逆転**: Repository インターフェースは Domain 層、実装は Infrastructure 層
- **テスタビリティ**: モック（mockery）で各レイヤーを独立テスト可能
