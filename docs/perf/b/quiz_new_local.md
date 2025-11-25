**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **200ms**
- Error rate < **1%**
- 対象: `POST /quizzes/new`（生成：CPU/DB 負荷）

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

- シナリオ: quiz_new（クイズ生成）
- ステージ: `1 VU × 1 iteration`（動作確認用）
- ThinkTime: 1s
- クイズパラメータ:
  - questionCount: 10
  - isSaveResult: false
  - isRegisteredWords: 0
  - correctRate: 1
  - partsOfSpeeches: [1]
- 閾値（thresholds）:
  - `http_req_duration{endpoint:quiz_new}: p(95)<200ms`
  - `http_req_failed{endpoint:quiz_new}: rate<0.01`
  - 注: 動作確認用のため、1 VU × 1 iteration のみ実行

**結果（要約）**

- SLO: **達成** ✅
- チェック成功率: **100%** (1/1)
- エラー率: **0.00%** (0/2)
- 実行時間: **約 1.1 秒**
- 注: このテストは動作確認用のため、1 VU × 1 iteration のみ実行

**エンドポイント別メトリクス**

| Endpoint          | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | SLO |
| ----------------- | -------: | -------: | -------: | ---------: | :-: |
| POST /quizzes/new |   **32** |   **45** |   **47** |  **0.00%** | ✅  |

注: このテストでは複数のエンドポイントを順次呼び出します（test-login、quiz_new）

**HTTP メトリクス**

- http_req_duration:
  - avg: **31.85ms**
  - min: **16.83ms**
  - med: **31.85ms**
  - max: **46.88ms**
  - p(90): **43.87ms**
  - p(95): **45.38ms**
- http_req_failed: **0.00%** (0/2)
- http_reqs: **2** (1.9 req/s)

**チェック結果**

- ✓ 2xx: **100%** (1/1)

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得
2. `POST /quizzes/new` - クイズを生成

**解釈 / ボトルネック**

- エラー率は **0.00%** で安定しています
- p95 latency は **45.38ms** で 200ms を大きく下回っており、SLO を達成しています
- 平均レイテンシは **31.85ms** で、クイズ生成処理としては非常に良好な性能です
- 最大レイテンシ（46.88ms）も 200ms を大きく下回っており、CPU/DB 負荷が高い処理にもかかわらず優秀な性能を示しています
- ただし、1 iteration のみの実行のため、統計的な信頼性は限定的です
- 本番相当の負荷テスト（PR: 10VU、Nightly: 30-50VU）では、より詳細な性能評価が必要です
- 特に、questionCount を増やしたり、isRegisteredWords を有効にした場合の性能変化も確認すべきです

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
export BASE_URL="http://localhost:8080"
export PROFILE="pr"
k6 run bench/k6/b/quiz_new.js
```
