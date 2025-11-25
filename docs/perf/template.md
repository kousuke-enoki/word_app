## SLO & Results

**SLO（目標）**

- p95 latency < **200ms**
- Error rate < **1%**
- 対象: `POST /api/v1/auth/login`, `GET /api/v1/words?q=`, `POST /api/v1/quizzes`

**計測メタデータ**

- **Date (JST)**: YYYY-MM-DD HH:MM:SS JST
- **Git SHA**: `abcdef1`
- **Go Version**: 1.23.x
- **k6 Version**: v0.50.x
- **DB Version**: PostgreSQL
- **DB Size**: words=260,000 / users=固定シード
- **Seed Type**: 固定シード

**テスト環境**

- Frontend: Vercel（東京近傍 CDN）
- Backend: API Gateway + Lambda (Go 1.23) **Provisioned Concurrency=1**
- DB: PostgreSQL **(行数: words=260,000 / users=固定シード)**
- リージョン: ap-northeast-1
- キャッシュ: なし（DB 直）

**ワークロード（k6）**

- シナリオ: login / words search / quiz generation
- ステージ: `1m → 10VU`, `2m → 30VU`, `1m → 0VU`
- ThinkTime: 1s
- 閾値（thresholds）:
  - `http_req_duration: p(95)<200ms`
  - `http_req_failed: rate<0.01`

**結果（要約）**

- SLO: **達成 / 未達成（どちらか）**
- 最大安定 RPS（30VU 時）: **XX req/s**
- コールドスタート: **除外 / 含む（どちらか）**（初回のみ +XXXms）

**エンドポイント別メトリクス**

| Endpoint         | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | 30VU 時 RPS |  SLO  |
| ---------------- | -------: | -------: | -------: | ---------: | ----------: | :---: |
| POST /auth/login |   **xx** |   **xx** |   **xx** |  **0.xx%** |     **x.x** | ✅/❌ |
| GET /words?q=... |   **xx** |   **xx** |   **xx** |  **0.xx%** |     **x.x** | ✅/❌ |
| POST /quizzes    |   **xx** |   **xx** |   **xx** |  **0.xx%** |     **x.x** | ✅/❌ |

**グラフ**

- RPS–Latency（search）：  
  ![search p95](search_p95.png)
- Error rate 推移（全体）：  
  ![error rate](error_rate.png)

**解釈 / ボトルネック**

- `GET /words` は **インデックス最適化** で p95 が **XXXms → YYYms (-ZZ%)**（詳細は [DB Perf](#)）
- `POST /quizzes` は **N+1 排除** で RPS 安定域が **+AA%**
- 残課題：コールドスタート影響（初回 +~300ms）→ **プロビジョンド維持**で回避

**再現手順**

```bash
# 前提: JMdict 26万件を導入済み（docs/data.md）
export BASE_URL="https://demo.taplex.app"
k6 run bench/k6/login.js   --out json=bench/k6/out/login.json
k6 run bench/k6/search.js  --out json=bench/k6/out/search.json
k6 run bench/k6/quiz.js    --out json=bench/k6/out/quiz.json
node scripts/print-k6-summary.mjs bench/k6/out/*.json
```
