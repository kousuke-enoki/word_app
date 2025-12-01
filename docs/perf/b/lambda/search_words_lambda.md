**Environment: Lambda (API Gateway + Lambda 256MB + RDS t4g.micro, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **1000ms**（Lambda環境向け：コールドスタートとネットワークレイテンシを考慮）
- Error rate < **1%**
- 対象: `GET /words?search=`（検索：読み多・索引効く）

**計測メタデータ**

- **Date (JST)**: YYYY-MM-DD HH:MM:SS JST
- **Git SHA**: `abcdef1`
- **Go Version**: 1.25.x
- **k6 Version**: v1.3.0
- **DB Version**: PostgreSQL（RDS t4g.micro）
- **DB Size**: words=要確認件数 / users=固定シード
- **Seed Type**: 固定シード

**テスト環境**

- Frontend: Vercel（東京近傍 CDN）
- Backend: API Gateway + Lambda (Go 1.25.x, 256MB, 30秒タイムアウト)
- DB: PostgreSQL（RDS t4g.micro）
- リージョン: ap-northeast-1
- キャッシュ: なし（DB 直）
- Provisioned Concurrency: なし（コールドスタート含む）

**ワークロード（k6）**

- シナリオ: search_words（単語検索）
- Executor: `constant-arrival-rate`
- Rate: **5.00 iterations/s**（PR プロファイル）または **10.00 iterations/s**（Nightly プロファイル）
- Duration: **2m0s**（PR）または **3m0s**（Nightly）
- VUs: 最大 **5**（PR）または **10**（Nightly）、preAllocated: **1**（Lambda向け：ウォームアップ用）
- ThinkTime: 1s
- 検索パラメータ: ランダムな検索クエリとソート条件（`randomSearchQuery()`, `randomSortBy()`）
- 閾値（thresholds）:
  - `http_req_duration{endpoint:search}: p(95)<1000ms`（Lambda向け）
  - `http_req_failed{endpoint:search}: rate<0.01`

**結果（要約）**

- SLO: **達成 / 未達成（どちらか）**
- チェック成功率: **XX%** (XXX/XXX)
- エラー率: **0.XX%** (X/XXX) - search エンドポイントのみ
- 実行時間: **約 X 分 XX 秒**
- イテレーション: **XXX** (X.XX iter/s)
- ドロップしたイテレーション: **XX** (X.XX iter/s) - VU 不足によりスキップ

**エンドポイント別メトリクス**

| Endpoint           | p50 (ms) | p95 (ms) | p99 (ms) | Error rate |      RPS | SLO |
| ------------------ | -------: | -------: | -------: | ---------: | -------: | :-: |
| GET /words?search= |   **xx** |   **xx** |   **xx** |  **0.xx%** | **x.xx** | ✅/❌ |

**HTTP メトリクス**

- http_req_duration (search):
  - avg: **XXX.XXms**
  - min: **XX.XXms**
  - med: **XXX.XXms**
  - max: **XXXX.XXms**
  - p(90): **XXX.XXms**
  - p(95): **XXX.XXms**
- http_req_failed (search): **0.XX%** (X/XXX)
- http_reqs: **XXX** (X.XX req/s)
- イテレーション: **XXX** (X.XX iter/s)

**チェック結果**

- ✓ 2xx: **XX%** (XXX/XXX)

**ネットワーク**

- 受信データ: **X.X MB** (XX kB/s)
- 送信データ: **XXX kB** (X.X kB/s)

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得（setup 関数）
2. `GET /words?search={ランダム}&sortBy={ランダム}` - 単語を検索（各イテレーション）

**解釈 / ボトルネック**

- Lambda環境での検索処理のパフォーマンス評価
- コールドスタートの影響: 初回リクエストで +XXXms の遅延
- API Gateway + Lambda + VPC + RDS のネットワークレイテンシを考慮
- 検索クエリのインデックス最適化の効果
- RDS t4g.micro での DB 接続プール（5接続）の影響
- ランダムな検索クエリとソート条件を使用しており、多様な検索パターンでの性能を評価

**再現手順**

```bash
# 前提: Lambda環境がデプロイ済み
export BASE_URL="https://xxxx.execute-api.ap-northeast-1.amazonaws.com/prod"
export PROFILE="pr"  # または "nightly"
k6 run bench/k6/b/search_words.js
```

