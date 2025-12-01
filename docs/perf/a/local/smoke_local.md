**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- Error rate < **1%**（全エンドポイント）
- 対象: 全公開・保護エンドポイントの動作確認
- 注: スモークテストは動作確認が目的のため、レイテンシ閾値（p95）は設定していません

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

- シナリオ: smoke test（全 API の動作確認）
- ステージ: `1 VU × 1 iteration`
- ThinkTime: なし
- 閾値（thresholds）:
  - `http_req_failed: rate==0`
  - 注: スモークテストは動作確認が目的のため、p95 レイテンシ閾値は設定していません

**結果（要約）**

- SLO: **達成** ✅
- チェック成功率: **100%** (9/9)
- エラー率: **0.00%** (0/10)
- 実行時間: **約 1.0 秒**

**エンドポイント別メトリクス**

| Endpoint                    | Status | Check Result | SLO |
| --------------------------- | ------ | :----------: | :-: |
| GET /health                 | 2xx    |      ✅      | ✅  |
| POST /users/auth/test-login | 2xx    |      ✅      | ✅  |
| GET /auth/check             | 2xx    |      ✅      | ✅  |
| GET /users/me               | 2xx    |      ✅      | ✅  |
| GET /users/my_page          | 2xx    |      ✅      | ✅  |
| GET /setting/user_config    | 2xx    |      ✅      | ✅  |
| GET /words?sortBy=name      | 2xx    |      ✅      | ✅  |
| GET /quizzes                | 2xx    |      ✅      | ✅  |
| GET /results                | 2xx    |      ✅      | ✅  |

**HTTP メトリクス**

- http_req_duration:
  - avg: **96.4ms**
  - min: **3.4ms**
  - med: **5.45ms**
  - max: **870.32ms**
  - p(90): **128.26ms**
  - p(95): **499.29ms**
- http_req_failed: **0.00%** (0/10)
- http_reqs: **10** (10.0 req/s)

**スキップされたエンドポイント（書き込み系）**

以下の書き込み系エンドポイントは `SMOKE_ENV=local` によりスキップされました：

- POST /setting/user_config
- POST /setting/root_config
- POST /words/register
- POST /words/memo
- POST /words/new
- PUT /words/1
- DELETE /words/1
- POST /words/bulk_tokenize
- POST /words/bulk_register
- POST /quizzes/new
- PUT /users/1
- DELETE /users/1

**解釈 / ボトルネック**

- 全エンドポイントが正常に動作し、エラーは発生していません
- p95 latency が 499.29ms とやや高いですが、これは smoke test（1 VU × 1 iteration）のため、初回リクエストやコールドスタートの影響が含まれている可能性があります
- 最大レイテンシ（870.32ms）は `/words?sortBy=name` などの DB クエリが重いエンドポイントの可能性があります

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
export BASE_URL="http://localhost:8080"
export SMOKE_ENV="local"
k6 run bench/k6/a/smoke_all.js
```
