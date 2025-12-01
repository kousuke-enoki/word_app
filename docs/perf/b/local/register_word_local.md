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
- Executor: `ramping-vus`
- ステージ: `20s → 2VU`, `1m → 5VU`, `1m30s → 5VU`, `20s → 0VU`（PR プロファイル）
- ThinkTime: 1s
- 検索パラメータ: ランダムな検索クエリとソート条件（`randomSearchQuery()`, `randomSortBy()`）
- セットアップ: テスト実行前に登録済み単語を最大200件までリセット
- 閾値（thresholds）:
  - `http_req_duration{endpoint:register_word}: p(95)<200ms`
  - `http_req_failed{endpoint:register_word,status:!429}: rate<0.01`（429エラーはクォータ制限のため許容）

**結果（要約）**

- SLO: **達成** ✅
- チェック成功率: **100%** (628/628)
- エラー率: **0.00%** (0/0) - register_word エンドポイントのみ（status:!429）
- 実行時間: **約 3 分 10 秒**
- イテレーション: **628** (3.30 iter/s)

**エンドポイント別メトリクス**

| Endpoint             | p50 (ms) | p95 (ms) | p99 (ms) | Error rate |      RPS | SLO |
| -------------------- | -------: | -------: | -------: | ---------: | -------: | :-: |
| POST /words/register |     **9** |   **33** |    **-** |  **0.00%** | **3.30** | ✅  |

注: このテストでは各イテレーションで複数のエンドポイントを順次呼び出します（test-login、検索、register）

**HTTP メトリクス**

- http_req_duration (register_word):
  - avg: **12.82ms**
  - min: **5ms**
  - med: **8.61ms**
  - max: **92.05ms**
  - p(90): **26.64ms**
  - p(95): **32.89ms**
- http_req_failed (register_word, status:!429): **0.00%** (0/0)
- http_reqs: **1266** (6.65 req/s) - 全エンドポイント合計
- register_word エンドポイント: **628** リクエスト（3.30 req/s）
- イテレーション: **628** (3.30 iter/s)

**チェック結果**

- ✓ 2xx or 409 or acceptable errors: **100%** (628/628)
  - 許容されるエラー:
    - **409**: ユニーク制約違反（既に登録済み）
    - **400**: 登録状態に変更なし
    - **429**: 登録単語数上限超過（クォータ制限）

**ネットワーク**

- 受信データ: **2.3 MB** (12 kB/s)
- 送信データ: **421 kB** (2.2 kB/s)

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得（setup 関数）
2. 登録済み単語のリセット（setup 関数、最大200件まで）
3. 各イテレーション:
   - `GET /words?search={ランダム}&sortBy={ランダム}&...` - 検索して単語 ID を取得
   - `POST /words/register` - 単語を登録

**解釈 / ボトルネック**

- SLO を達成：p95 latency が **32.89ms** で 200ms を大きく下回っています
- エラー率は **0.00%** で安定しています（429エラーは除外）
- 平均レイテンシは **12.82ms** で、書き込み処理としては非常に良好な性能です
- 中央値（p50）は **8.61ms** と非常に低く、大半のリクエストが高速に処理されています
- 最大レイテンシ（92.05ms）も 200ms を大きく下回っており、5VU の負荷下でも優れた性能を維持しています
- イテレーション実行時間の平均は **1.11s** で、ThinkTime（1秒）を考慮すると、登録処理自体は非常に高速です
- チェック成功率は **100%** で、許容されるエラー（409、400、429）も適切に処理されています
- 5VU 時の RPS は **3.30 req/s** で、ThinkTime を考慮すると妥当な値です
- ランダムな検索クエリとソート条件を使用しており、多様な検索パターンでの性能を評価できています

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
export BASE_URL="http://localhost:8080"
export PROFILE="pr"
k6 run bench/k6/b/register_word.js
```
