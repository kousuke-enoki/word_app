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
- **Profile**: **nightly**（最大10 VU）

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
  - Nightly: `30s → 1VU`（ウォームアップ）, `20s → 3VU`, `1m → 10VU`, `2m → 10VU`, `20s → 0VU`
- ThinkTime: 1s
- 閾値（thresholds）:
  - `http_req_duration{endpoint:sign_in}: p(95)<1000ms`（Lambda 向け）
  - `http_req_failed{endpoint:sign_in}: rate<0.01`

**結果（要約）**

- SLO: **達成** ✅
- 最大安定 RPS（10VU 時）: **4.28 req/s**
- コールドスタート: **含む**（ウォームアップ後も初回リクエストで +XXXms）
- チェック成功率: **100%** (2164/2164)
- エラー率: **0.00%** (0/1082) - sign_in エンドポイントのみ
- 実行時間: **約 4 分 12 秒**
- イテレーション: **1082** (4.28 iter/s)

**エンドポイント別メトリクス**

| Endpoint            | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | 10VU 時 RPS | SLO |
| ------------------- | -------: | -------: | -------: | ---------: | ----------: | :-: |
| POST /users/sign_in | **552.81** | **636.44** | **XXX** | **0.00%** | **4.28** | ✅ |

**HTTP メトリクス**

- http_req_duration (sign_in):
  - avg: **567.73ms**
  - min: **502.44ms**
  - med: **552.81ms**
  - max: **1.49s**
  - p(90): **615.87ms**
  - p(95): **636.44ms**
- http_req_failed (sign_in): **0.00%** (0/1082)
- http_reqs: **1083** (4.28 req/s)
- イテレーション: **1082** (4.28 iter/s)

**チェック結果**

- ✓ 2xx: **100%** (1082/1082)
- ✓ has token: **100%** (1082/1082)

**ネットワーク**

- 受信データ: **770 kB** (3.0 kB/s)
- 送信データ: **140 kB** (555 B/s)

**実行フロー**

1. `POST /users/sign_up` - テストユーザー作成（setup 関数、リトライ付き）
2. `POST /users/sign_in` - 認証（各イテレーション）

**解釈 / ボトルネック**

- Lambda 環境での認証処理のパフォーマンス評価（Nightly プロファイル: 最大10 VU）
- **SLO達成**: p95=636.44ms < 1000ms、エラー率0.00% < 1%
- **安定性**: 全イテレーション完了（1082件）、チェック成功率100%
- **レスポンス時間**: 平均567.73ms、中央値552.81ms、p95=636.44ms
  - Lambda 環境としてはやや高めだが、API Gateway + Lambda + VPC + RDS のネットワークレイテンシを考慮すると許容範囲
  - 最大値1.49秒は許容範囲内（コールドスタートやネットワーク遅延の影響）
  - 中央値と平均が近い（552.81ms vs 567.73ms）→ 分布が比較的安定
- **RPS**: 10VUで4.28 req/s（レスポンス時間平均567.73ms + ThinkTime 1秒 ≈ 1.57秒/イテレーション）
- コールドスタートの影響: ウォームアップフェーズ（30 秒、1 VU）により、コールドスタートの影響を軽減
- API Gateway + Lambda + VPC + RDS のネットワークレイテンシを考慮
- 256MB メモリでの処理性能は十分
- RDS t4g.micro での DB 接続プール（5 接続）の影響
- **推奨対応**:
  - Provisioned Concurrency の検討（コールドスタート影響の軽減）
  - より長いピーク負荷期間でのテスト（現在は2分間）
  - p99 の記録（現在はp95まで）
  - 複数回実行による再現性の確認（2回目の実行で外れ値が解消され、より安定した結果を確認）

**再現手順**

```bash
# 前提: Lambda環境がデプロイ済み
export BASE_URL="https://xxxx.execute-api.ap-northeast-1.amazonaws.com/prod"
export PROFILE="nightly"
export TEST_EMAIL="demo@example.com"
export TEST_PASSWORD="K6passw0rd!"
k6 run bench/k6/b/sign_in.js
```

