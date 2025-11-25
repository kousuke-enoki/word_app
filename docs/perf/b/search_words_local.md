**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **200ms**
- Error rate < **1%**
- 対象: `GET /words?q=`（検索：読み多・索引効く）

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
- ステージ: `1 VU × 1 iteration`（動作確認用）
- ThinkTime: 1s
- 検索パラメータ: `SEARCH_Q=test`, `SEARCH_SORT=name`
- 閾値（thresholds）:
  - `http_req_duration{endpoint:search}: p(95)<200ms`
  - `http_req_failed{endpoint:search}: rate<0.01`
  - 注: 動作確認用のため、1 VU × 1 iteration のみ実行

**結果（要約）**

- SLO: **達成** ✅
- チェック成功率: **100%** (1/1)
- エラー率: **0.00%** (0/2)
- 実行時間: **約 1.0 秒**
- 注: このテストは動作確認用のため、1 VU × 1 iteration のみ実行

**エンドポイント別メトリクス**

| Endpoint          | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | SLO |
| ----------------- | -------: | -------: | -------: | ---------: | :-: |
| GET /words?q=test |   **16** |   **18** |   **18** |  **0.00%** | ✅  |

注: このテストでは複数のエンドポイントを順次呼び出します（test-login、search）

**HTTP メトリクス**

- http_req_duration:
  - avg: **15.89ms**
  - min: **14.1ms**
  - med: **15.89ms**
  - max: **17.68ms**
  - p(90): **17.32ms**
  - p(95): **17.5ms**
- http_req_failed: **0.00%** (0/2)
- http_reqs: **2** (1.9 req/s)

**チェック結果**

- ✓ 2xx: **100%** (1/1)

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得
2. `GET /words?q=test&sortBy=name` - 単語を検索

**解釈 / ボトルネック**

- エラー率は **0.00%** で安定しています
- p95 latency は **17.5ms** で 200ms を大きく下回っており、SLO を達成しています
- 平均レイテンシは **15.89ms** で、検索処理としては非常に良好な性能です
- 最大レイテンシ（17.68ms）も 200ms を大きく下回っており、インデックスが効いていることを示しています
- 検索クエリ（`q=test`）は一般的な検索語であり、インデックス最適化の効果が確認できます
- ただし、1 iteration のみの実行のため、統計的な信頼性は限定的です
- 本番相当の負荷テスト（PR: 10VU、Nightly: 30-50VU）では、より詳細な性能評価が必要です
- 特に、大量のデータ（26 万件）に対する検索性能や、異なる検索クエリでの性能変化も確認すべきです

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
export BASE_URL="http://localhost:8080"
export PROFILE="pr"
export SEARCH_Q="test"
export SEARCH_SORT="name"
k6 run bench/k6/b/search_words.js
```
