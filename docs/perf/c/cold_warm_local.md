**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **200ms**（warm のみ評価）
- Error rate < **1%**（warm のみ評価）
- 対象: `GET /words?q=`（検索：コールドスタート vs ウォームスタート比較）

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

- シナリオ: cold/warm 比較（単語検索）
  - **cold_once**: 1 VU × 1 iteration（コールドスタート測定、記録のみ）
  - **warm_ramp**: 60 秒後に開始、PR プロファイル（`30s → 5VU`, `2m → 10VU`, `30s → 0VU`）
- ThinkTime: 1s
- 検索パラメータ: `SEARCH_Q=test`, `SEARCH_SORT=name`
- 閾値（thresholds）:
  - `http_req_duration{endpoint:search,phase:warm}: p(95)<200ms`
  - `http_req_failed{endpoint:search,phase:warm}: rate<0.01`
  - 注: cold は記録のみ（閾値評価対象外）

**結果（要約）**

- SLO: **達成** ✅
- 最大安定 RPS（10VU 時）: **4.5 req/s**
- チェック成功率: **100%** (1070/1070)
- エラー率: **0.00%** (0/1071)
- 実行時間: **約 4 分**

**エンドポイント別メトリクス**

| Endpoint          | Phase | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | 10VU 時 RPS | SLO |
| ----------------- | ----- | -------: | -------: | -------: | ---------: | ----------: | :-: |
| GET /words?q=test | cold  |   **15** |  **165** |  **165** |  **0.00%** |     **0.5** |  -  |
| GET /words?q=test | warm  |   **15** |   **20** |   **33** |  **0.00%** |     **4.5** | ✅  |

注: cold は記録のみで、SLO 評価対象外です

**HTTP メトリクス**

- **Cold Start（cold_once シナリオ）**:

  - http_req_duration: max=**164.84ms**（全体統計から推定）
  - http_req_failed: **0.00%**
  - 注: コールドスタート時の初回リクエストのレイテンシを測定（1 iteration のみのため統計的信頼性は限定的）

- **Warm Start（warm_ramp シナリオ）**:

  - http_req_duration:
    - avg: **14.04ms**
    - min: **4.33ms**
    - med: **15.43ms**
    - max: **32.9ms**
    - p(90): **18.27ms**
    - p(95): **19.58ms**
  - http_req_failed: **0.00%** (0/1069)
  - http_reqs: **1069** (4.5 req/s)
  - 注: ウォーム状態での負荷テスト結果（10VU 時）

- **全体（cold + warm）**:
  - http_req_duration:
    - avg: **14.17ms**
    - min: **4.33ms**
    - med: **15.43ms**
    - max: **164.84ms**
    - p(90): **18.29ms**
    - p(95): **19.62ms**
  - http_req_failed: **0.00%** (0/1071)
  - http_reqs: **1071** (4.5 req/s)
  - イテレーション: **1070** (4.4 iter/s)

**コールドスタート影響**

- コールドスタート時の最大レイテンシ: **164.84ms**（全体統計から推定）
- ウォーム状態の p95 レイテンシ: **19.58ms**
- コールドスタート時の追加レイテンシ: **約 +145ms**（推定）
- 注: cold phase は 1 iteration のみのため、統計的な信頼性は限定的です

**チェック結果**

- ✓ 2xx: **100%** (1070/1070)

**解釈 / ボトルネック**

- **SLO 達成**: warm phase の p95 latency が **19.58ms** で 200ms を大きく下回っており、SLO を達成しています
- **エラー率**: **0.00%** で安定しています
- **コールドスタート影響**: コールドスタート時の最大レイテンシは **164.84ms** で、ウォーム状態の p95（19.58ms）と比較して約 **+145ms** の追加レイテンシが発生しています
- **ウォーム状態の性能**: 平均レイテンシ **14.04ms**、p95 レイテンシ **19.58ms** と非常に良好な性能を示しています
- **ローカル環境での測定**: ローカル環境では Lambda のコールドスタートは発生しませんが、アプリケーション初期化や DB 接続プールの初期化などの影響を測定できます
- **本番環境での考慮**: 本番環境（Lambda）では、Provisioned Concurrency を有効にすることでコールドスタートを回避できます
- **統計的な注意**: cold phase は 1 iteration のみのため、統計的な信頼性は限定的です。より正確な測定には複数回の実行が必要です

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
# 事前にアイドル状態を作るため、60秒以上アクセスがない状態にする

export BASE_URL="http://localhost:8080"
export SEARCH_Q="test"
export SEARCH_SORT="name"
npm --prefix bench run k6:c:cold

# 結果は bench/k6/out/cold_warm.json に出力されます
```

**注意事項**

- このテストは 2 つのシナリオを並行実行します：
  1. `cold_once`: 0 秒時に開始、1 VU × 1 iteration
  2. `warm_ramp`: 60 秒後に開始、PR プロファイルで負荷テスト
- コールドスタートの測定には、事前に一定時間（60 秒以上）アクセスがない状態を作る必要があります
- 本番環境（Lambda）での実行時は、Provisioned Concurrency の有無で結果が大きく異なります
