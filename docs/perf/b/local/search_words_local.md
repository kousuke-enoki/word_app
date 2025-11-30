**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **200ms**
- Error rate < **1%**
- 対象: `GET /words?search=`（検索：読み多・索引効く）

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

- シナリオ: search_words（単語検索）
- Executor: `constant-arrival-rate`
- Rate: **5.00 iterations/s**（PR プロファイル）
- Duration: **2m0s**
- VUs: 最大 **5**（preAllocated: 3）
- ThinkTime: 1s
- 検索パラメータ: ランダムな検索クエリとソート条件（`randomSearchQuery()`, `randomSortBy()`）
- 閾値（thresholds）:
  - `http_req_duration{endpoint:search}: p(95)<200ms`
  - `http_req_failed{endpoint:search}: rate<0.01`

**結果（要約）**

- SLO: **達成** ✅
- チェック成功率: **100%** (492/492)
- エラー率: **0.00%** (0/492) - search エンドポイントのみ
- 実行時間: **約 2 分 1 秒**
- イテレーション: **492** (4.07 iter/s)
- ドロップしたイテレーション: **108** (0.89 iter/s) - VU 不足によりスキップ

**エンドポイント別メトリクス**

| Endpoint           | p50 (ms) | p95 (ms) | p99 (ms) | Error rate |      RPS | SLO |
| ------------------ | -------: | -------: | -------: | ---------: | -------: | :-: |
| GET /words?search= |   **88** |  **148** |    **-** |  **0.00%** | **4.08** | ✅  |

**HTTP メトリクス**

- http_req_duration (search):
  - avg: **92.37ms**
  - min: **50.22ms**
  - med: **87.57ms**
  - max: **664.03ms**
  - p(90): **111.36ms**
  - p(95): **147.92ms**
- http_req_failed (search): **0.00%** (0/492)
- http_reqs: **493** (4.08 req/s)
- イテレーション: **492** (4.07 iter/s)

**チェック結果**

- ✓ 2xx: **100%** (492/492)

**ネットワーク**

- 受信データ: **2.1 MB** (18 kB/s)
- 送信データ: **149 kB** (1.2 kB/s)

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得（setup 関数）
2. `GET /words?search={ランダム}&sortBy={ランダム}` - 単語を検索（各イテレーション）

**解釈 / ボトルネック**

- SLO を達成：p95 latency が **147.92ms** で 200ms を下回っています
- エラー率は **0.00%** で安定しています
- 平均レイテンシは **92.37ms** で、検索処理としては良好な性能です
- 中央値（p50）は **87.57ms** と低く、大半のリクエストが高速に処理されています
- 最大レイテンシ（664.03ms）は散発的に発生していますが、p95 が 200ms を下回っているため許容範囲内です
- 目標レート 5.00 iter/s に対して、実際のイテレーション数は **4.07 iter/s** でした
- **108 件のイテレーションがドロップ**されました（VU 不足により）。これは maxVUs=5 では目標レートを維持するには不十分だったことを示しています
- VU は最小 3、最大 5 で動作し、平均 4 VU で実行されました
- ランダムな検索クエリとソート条件を使用しており、多様な検索パターンでの性能を評価できています

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
export BASE_URL="http://localhost:8080"
export PROFILE="pr"
k6 run bench/k6/b/search_words.js
```
