**Environment: Lambda (API Gateway + Lambda 256MB + RDS t4g.micro, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **1000ms**（Lambda 環境向け：コールドスタートとネットワークレイテンシを考慮）
- Error rate < **1%**
- 対象: `POST /users/sign_in`（認証：外部 I/O ほぼ無し、基準線）

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
- Backend: API Gateway + Lambda (Go 1.25.x, 256MB, 30 秒タイムアウト)
- DB: PostgreSQL（RDS t4g.micro）
- リージョン: ap-northeast-1
- キャッシュ: なし（DB 直）
- Provisioned Concurrency: なし（コールドスタート含む）

**ワークロード（k6）**

- シナリオ: sign_in（認証）
- Executor: `ramping-vus`
- ステージ（Lambda 向け: ウォームアップ付き）:
  - PR: `30s → 1VU`（ウォームアップ）, `20s → 2VU`, `1m → 5VU`, `1m30s → 5VU`, `20s → 0VU`
  - Nightly: `30s → 1VU`（ウォームアップ）, `20s → 3VU`, `1m → 10VU`, `2m → 10VU`, `20s → 0VU`
- ThinkTime: 1s
- 閾値（thresholds）:
  - `http_req_duration{endpoint:sign_in}: p(95)<1000ms`（Lambda 向け）
  - `http_req_failed{endpoint:sign_in}: rate<0.01`

**結果（要約）**

- SLO: **達成 / 未達成（どちらか）**
- 最大安定 RPS（5VU/10VU 時）: **XX req/s**
- コールドスタート: **含む**（ウォームアップ後も初回リクエストで +XXXms）
- チェック成功率: **XX%** (XXX/XXX)
- エラー率: **0.XX%** (X/XXX) - sign_in エンドポイントのみ
- 実行時間: **約 X 分 XX 秒**

**エンドポイント別メトリクス**

| Endpoint            | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | 5VU/10VU 時 RPS |  SLO  |
| ------------------- | -------: | -------: | -------: | ---------: | --------------: | :---: |
| POST /users/sign_in |   **xx** |   **xx** |   **xx** |  **0.xx%** |         **x.x** | ✅/❌ |

**HTTP メトリクス**

- http_req_duration (sign_in):
  - avg: **XX.XXms**
  - min: **XX.XXms**
  - med: **XX.XXms**
  - max: **XXXX.XXms**
  - p(90): **XXX.XXms**
  - p(95): **XXX.XXms**
- http_req_failed (sign_in): **0.XX%** (X/XXX)
- http_reqs: **XXX** (X.XX req/s)
- イテレーション: **XXX** (X.XX iter/s)

**チェック結果**

- ✓ 2xx: **XX%** (XXX/XXX)
- ✓ has token: **XX%** (XXX/XXX)

**ネットワーク**

- 受信データ: **X.X MB** (XX kB/s)
- 送信データ: **XXX kB** (X.X kB/s)

**実行フロー**

1. `POST /users/sign_up` - テストユーザー作成（setup 関数、リトライ付き）
2. `POST /users/sign_in` - 認証（各イテレーション）

**解釈 / ボトルネック**

- Lambda 環境での認証処理のパフォーマンス評価
- コールドスタートの影響: 初回リクエストで +XXXms の遅延
- API Gateway + Lambda + VPC + RDS のネットワークレイテンシを考慮
- ウォームアップフェーズ（30 秒、1 VU）により、コールドスタートの影響を軽減
- 256MB メモリでの処理性能
- RDS t4g.micro での DB 接続プール（5 接続）の影響

**再現手順**

```bash
# 前提: Lambda環境がデプロイ済み
export BASE_URL="https://xxxx.execute-api.ap-northeast-1.amazonaws.com/prod"
export PROFILE="pr"  # または "nightly"
export TEST_EMAIL="demo@example.com"
export TEST_PASSWORD="K6passw0rd!"
k6 run bench/k6/b/sign_in.js
```
