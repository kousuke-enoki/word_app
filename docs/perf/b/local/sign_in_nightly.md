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
- Executor: `ramping-vus`
- ステージ: `20s → 3VU`, `1m → 10VU`, `2m → 10VU`, `20s → 0VU`（Nightly プロファイル）
- ThinkTime: 1s
- 閾値（thresholds）:
  - `http_req_duration{endpoint:sign_in}: p(95)<200ms`
  - `http_req_failed{endpoint:sign_in}: rate<0.01`

**結果（要約）**

- SLO: **達成** ✅
- 最大安定 RPS（10VU 時）: **7.08 req/s**
- チェック成功率: **100%** (3120/3120)
- エラー率: **0.00%** (0/1560) - sign_in エンドポイントのみ
- 実行時間: **約 3 分 40 秒**

**エンドポイント別メトリクス**

| Endpoint            | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | 10VU 時 RPS | SLO |
| ------------------- | -------: | -------: | -------: | ---------: | ----------: | :-: |
| POST /users/sign_in |   **81** |  **105** |    **-** |  **0.00%** |    **7.08** | ✅  |

**HTTP メトリクス**

- http_req_duration (sign_in):
  - avg: **84.77ms**
  - min: **64.22ms**
  - med: **81.49ms**
  - max: **472.04ms**
  - p(90): **93.43ms**
  - p(95): **105.48ms**
- http_req_failed (sign_in): **0.00%** (0/1560)
- http_reqs: **1561** (7.08 req/s)
- イテレーション: **1560** (7.07 iter/s)

**チェック結果**

- ✓ 2xx: **100%** (3120/3120)
- ✓ has token: **100%** (3120/3120)

**ネットワーク**

- 受信データ: **573 kB** (2.6 kB/s)
- 送信データ: **292 kB** (1.3 kB/s)

**解釈 / ボトルネック**

- SLO を達成：p95 latency が **105.48ms** で 200ms を大きく下回っています
- エラー率は **0.00%** で安定しています
- 平均レイテンシは **84.77ms** で、認証処理としては良好な性能です
- 中央値（p50）は **81.49ms** と低く、大半のリクエストが高速に処理されています
- 最大レイテンシ（472.04ms）は 10VU 時の負荷下でも散発的に発生していますが、p95 が 200ms を下回っているため許容範囲内です
- 10VU 時の RPS は **7.08 req/s** で、ThinkTime（1 秒）を考慮すると妥当な値です
- PR プロファイル（5VU）と比較して、VU 数が 2 倍になったことで RPS も約 2 倍（3.39 → 7.08）になっていますが、レイテンシは同程度を維持しており、スケーラビリティが良好です

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
export BASE_URL="http://localhost:8080"
export PROFILE="nightly"
export TEST_EMAIL="demo@example.com"
export TEST_PASSWORD="Secret-k6"
k6 run bench/k6/b/sign_in.js
```

