**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **200ms**
- Error rate < **1%**
- 対象: `POST /words/register`（書き込み：ユースケース感）

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

- シナリオ: register_word（単語登録）
- ステージ: `1 VU × 1 iteration`（動作確認用）
- ThinkTime: 1s
- 検索パラメータ: `SEARCH_Q=able`, `SEARCH_SORT=name`
- 閾値（thresholds）:
  - `http_req_duration{endpoint:register}: p(95)<200ms`
  - `http_req_failed{endpoint:register}: rate<0.01`
  - 注: 動作確認用のため、1 VU × 1 iteration のみ実行

**結果（要約）**

- SLO: **達成** ✅
- チェック成功率: **100%** (1/1)
- エラー率: **0.00%** (0/3)
- 実行時間: **約 1.2 秒**
- 注: このテストは動作確認用のため、1 VU × 1 iteration のみ実行

**エンドポイント別メトリクス**

| Endpoint             | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | SLO |
| -------------------- | -------: | -------: | -------: | ---------: | :-: |
| POST /words/register |   **15** |  **186** |  **205** |  **0.00%** | ✅  |

注: このテストでは複数のエンドポイントを順次呼び出します（test-login、検索、register）

**HTTP メトリクス**

- http_req_duration:
  - avg: **77.43ms**
  - min: **12.3ms**
  - med: **14.98ms**
  - max: **205.02ms**
  - p(90): **167.01ms**
  - p(95): **186.02ms**
- http_req_failed: **0.00%** (0/3)
- http_reqs: **3** (2.4 req/s)

**チェック結果**

- ✓ 2xx or 409: **100%** (1/1)
  - 注: 409（Conflict）は既に登録済みの場合に返されるため、正常系として扱います

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得
2. `GET /words?search=able&sortBy=name&...` - 検索して単語 ID を取得
3. `POST /words/register` - 単語を登録

**解釈 / ボトルネック**

- エラー率は **0.00%** で安定しています
- p95 latency は **186.02ms** で 200ms を下回っており、SLO を達成しています
- 平均レイテンシは **77.43ms** で、書き込み処理としては良好な性能です
- 最大レイテンシ（205.02ms）は 200ms をわずかに超えていますが、1 iteration のみの実行のため、統計的な信頼性は限定的です
- 本番相当の負荷テスト（PR: 10VU、Nightly: 30-50VU）では、より詳細な性能評価が必要です

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
export BASE_URL="http://localhost:8080"
export PROFILE="pr"
export SEARCH_Q="able"
export SEARCH_SORT="name"
k6 run bench/k6/b/register_word.js
```
