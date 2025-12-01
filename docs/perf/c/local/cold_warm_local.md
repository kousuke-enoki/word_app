**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **200ms**（warm のみ評価）
- Error rate < **1%**（warm のみ評価）
- 対象: `GET /words?search=`（検索：コールドスタート vs ウォームスタート比較）

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
  - **warm_ramp**: 60 秒後に開始、PR プロファイル（`20s → 2VU`, `1m30s → 5VU`, `20s → 0VU`）
- ThinkTime: 1s
- 検索パラメータ: `SEARCH_Q=test`（デフォルト）、`SEARCH_SORT=name`（デフォルト）
- 閾値（thresholds）:
  - `http_req_duration{endpoint:search,phase:warm}: p(95)<200ms`
  - `http_req_failed{endpoint:search,phase:warm}: rate<0.01`
  - 注: cold は記録のみ（閾値評価対象外）

**結果（要約）**

- SLO: **達成** ✅
- 最大安定 RPS（5VU 時）: **1.69 req/s**
- チェック成功率: **100%** (321/321)
- エラー率: **0.00%** (0/322) - warm phase のみ
- 実行時間: **約 3 分 10 秒**

**エンドポイント別メトリクス**

| Endpoint           | Phase | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | 5VU 時 RPS | SLO |
| ------------------ | ----- | -------: | -------: | -------: | ---------: | ---------: | :-: |
| GET /words?search= | cold  |    **-** |    **-** |    **-** |  **0.00%** |   **0.01** |  -  |
| GET /words?search= | warm  |    **94** |  **131** |    **-** |  **0.00%** |   **1.69** | ✅  |

注: cold は記録のみで、SLO 評価対象外です

**HTTP メトリクス**

- **Cold Start（cold_once シナリオ）**:

  - http_req_duration: 1 iteration のみのため統計なし
  - http_req_failed: **0.00%**
  - 注: コールドスタート時の初回リクエストのレイテンシを測定（1 iteration のみのため統計的信頼性は限定的）

- **Warm Start（warm_ramp シナリオ）**:

  - http_req_duration (warm):
    - avg: **100.83ms**
    - min: **74.95ms**
    - med: **94.4ms**
    - max: **543.93ms**
    - p(90): **121.22ms**
    - p(95): **131.19ms**
  - http_req_failed (warm): **0.00%** (0/320)
  - http_reqs: **322** (1.69 req/s) - 全エンドポイント合計
  - warm phase のリクエスト: **320** リクエスト
  - 注: ウォーム状態での負荷テスト結果（最大5VU 時）

- **全体（cold + warm）**:
  - http_req_duration:
    - avg: **104.1ms**
    - min: **74.95ms**
    - med: **94.44ms**
    - max: **685.44ms**
    - p(90): **121.4ms**
    - p(95): **131.84ms**
  - http_req_failed: **0.00%** (0/322)
  - http_reqs: **322** (1.69 req/s)
  - イテレーション: **321** (1.68 iter/s)

**コールドスタート影響**

- コールドスタート時: 1 iteration のみのため、詳細な統計は取得できませんでした
- ウォーム状態の p95 レイテンシ: **131.19ms**
- 注: cold phase は 1 iteration のみのため、統計的な信頼性は限定的です

**チェック結果**

- ✓ 2xx: **100%** (321/321)

**ネットワーク**

- 受信データ: **1.0 MB** (5.3 kB/s)
- 送信データ: **95 kB** (499 B/s)

**解釈 / ボトルネック**

- **SLO 達成**: warm phase の p95 latency が **131.19ms** で 200ms を下回っており、SLO を達成しています
- **エラー率**: **0.00%** で安定しています
- **ウォーム状態の性能**: 平均レイテンシ **100.83ms**、p95 レイテンシ **131.19ms** で良好な性能を示しています
- **中央値（p50）**: **94.4ms** と低く、大半のリクエストが高速に処理されています
- **最大レイテンシ**: **543.93ms** は散発的に発生していますが、p95 が 200ms を下回っているため許容範囲内です
- **コールドスタート測定**: cold phase は 1 iteration のみのため、詳細な統計は取得できませんでした
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
  1. `cold_once`: 0 秒時に開始、1 VU × 1 iteration（約 1.6 秒で完了）
  2. `warm_ramp`: 60 秒後に開始、PR プロファイル（最大5VU）で負荷テスト（2分10秒間）
- コールドスタートの測定には、事前に一定時間（60 秒以上）アクセスがない状態を作る必要があります
- 本番環境（Lambda）での実行時は、Provisioned Concurrency の有無で結果が大きく異なります
