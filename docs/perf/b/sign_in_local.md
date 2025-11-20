**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **200ms**
- Error rate < **1%**
- 対象: `POST /users/sign_in`

**計測メタデータ**

- **Date (JST)**: 2025-11-25 16:59:11 JST
- **Git SHA**: `e672e06`
- **Go Version**: 1.25.4
- **k6 Version**: v1.3.0
- **DB Version**: PostgreSQL（ローカル）
- **DB Size**: words=要確認件数 / users=固定シード
- **Seed Type**: 固定シード

**テスト環境**

- Frontend: ローカル（未使用）
- Backend: ローカル（Go 1.25.4）
- DB: PostgreSQL（ローカル）
- リージョン: ローカル
- キャッシュ: なし（DB 直）

**ワークロード（k6）**

- シナリオ: sign_in（認証）
- ステージ: `30s → 5VU`, `2m → 10VU`, `30s → 0VU`（PR プロファイル）
- ThinkTime: 1s
- 閾値（thresholds）:
  - `http_req_duration{endpoint:sign_in}: p(95)<200ms`
  - `http_req_failed{endpoint:sign_in}: rate<0.01`

**結果（要約）**

- SLO: **達成** ✅
- 最大安定 RPS（10VU 時）: **5.4 req/s**
- チェック成功率: **100%** (1940/1940)
- エラー率: **0.00%** (0/970) - sign_in エンドポイントのみ
- 実行時間: **約 3 分**

**エンドポイント別メトリクス**

| Endpoint            | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | 10VU 時 RPS | SLO |
| ------------------- | -------: | -------: | -------: | ---------: | ----------: | :-: |
| POST /users/sign_in |  **114** |  **165** |  **240** |  **0.00%** |     **5.4** | ✅  |

**HTTP メトリクス**

- http_req_duration (sign_in):
  - avg: **117.04ms**
  - min: **69.27ms**
  - med: **114.25ms**
  - max: **239.64ms**
  - p(90): **149.71ms**
  - p(95): **164.68ms**
- http_req_failed (sign_in): **0.00%** (0/970)
- http_reqs: **971** (5.4 req/s)
- イテレーション: **970** (5.4 iter/s)

**チェック結果**

- ✓ 2xx: **100%** (1940/1940)
- ✓ has token: **100%** (1940/1940)

**解釈 / ボトルネック**

- SLO を達成：p95 latency が **164.68ms** で 200ms を下回っています
- エラー率は **0.00%** で安定しています
- 平均レイテンシは **117.04ms** で、認証処理としては良好な性能です
- 最大レイテンシ（239.64ms）は 10VU 時の負荷下でも 200ms をわずかに超える程度で、許容範囲内です
- 10VU 時の RPS は **5.4 req/s** で、ThinkTime（1 秒）を考慮すると妥当な値です

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
export BASE_URL="http://localhost:8080"
export PROFILE="pr"
export TEST_EMAIL="demo@example.com"
export TEST_PASSWORD="Secret-k6"
k6 run bench/k6/b/sign_in.js
```
