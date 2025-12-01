**Environment: Lambda (API Gateway + Lambda 256MB + RDS t4g.micro, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **1000ms**（Lambda環境向け：コールドスタートとネットワークレイテンシを考慮）
- Error rate < **1%**（429エラーはクォータ制限のため許容）
- 対象: `POST /quizzes/new`（生成：CPU/DB 負荷）

**計測メタデータ**

- **Date (JST)**: YYYY-MM-DD HH:MM:SS JST
- **Git SHA**: `abcdef1`
- **Go Version**: 1.25.x
- **k6 Version**: v1.3.0
- **DB Version**: PostgreSQL（RDS t4g.micro）
- **DB Size**: words=要確認件数 / users=固定シード
- **Seed Type**: 固定シード

**テスト環境**

- Frontend: Vercel（東京近傍 CDN）
- Backend: API Gateway + Lambda (Go 1.25.x, 256MB, 30秒タイムアウト)
- DB: PostgreSQL（RDS t4g.micro）
- リージョン: ap-northeast-1
- キャッシュ: なし（DB 直）
- Provisioned Concurrency: なし（コールドスタート含む）

**ワークロード（k6）**

- シナリオ: quiz_new（クイズ生成）
- Executor: `ramping-vus`
- ステージ（Lambda向け: ウォームアップ付き）:
  - PR: `30s → 1VU`（ウォームアップ）, `20s → 2VU`, `1m → 5VU`, `1m30s → 5VU`, `20s → 0VU`
  - Nightly: `30s → 1VU`（ウォームアップ）, `20s → 3VU`, `1m → 10VU`, `2m → 10VU`, `20s → 0VU`
- ThinkTime: **5-10秒**（ランダム、クイズ生成は重い処理のため）
- クイズパラメータ: ランダムなクイズパラメータ（`randomQuizParams()`）
  - questionCount: 10 または 20（ランダム）
  - isSaveResult: false
  - isRegisteredWords: 0（全単語）
  - correctRate: 1
  - partsOfSpeeches: [1]
  - isIdioms: 0 または 1（ランダム）
  - isSpecialCharacters: 0 または 1（ランダム）
- 閾値（thresholds）:
  - `http_req_duration{endpoint:quiz_new}: p(95)<1000ms`（Lambda向け）
  - `http_req_failed{endpoint:quiz_new,status:!429}: rate<0.01`（429エラーはクォータ制限のため許容）

**結果（要約）**

- SLO: **達成 / 未達成（どちらか）**
- チェック成功率: **XX%** (XXX/XXX)
- エラー率: **XX.XX%** (XX/XXX) - quiz_new エンドポイントのみ
  - 主なエラー: 500 internal error、409 another quiz is running、429（クォータ制限）
- 実行時間: **約 X 分 XX 秒**
- イテレーション: **XX** (X.XX iter/s)
- 成功したイテレーション: **XX 件**（2xx ステータスのみ）

**エンドポイント別メトリクス**

| Endpoint          | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | RPS | SLO |
| ----------------- | -------: | -------: | -------: | ---------: | --: | :-: |
| POST /quizzes/new |  **xxx** |  **xxx** |   **xx** | **xx.xx%** | **x.xx** | ✅/❌ |

注: このテストでは各イテレーションで複数のエンドポイントを順次呼び出します（test-login、quiz_new）

**HTTP メトリクス**

- http_req_duration (quiz_new):
  - avg: **X.XXs**
  - min: **XXX.XXms**
  - med: **XXX.XXms**
  - max: **XX.XXs**
  - p(90): **X.XXs**
  - p(95): **X.XXs**
- http_req_failed (quiz_new): **XX.XX%** (XX/XXX)
  - 429エラー: **X 件**（クォータ制限によるエラー）
  - 主なエラー: 500 internal error、409 another quiz is running
- http_reqs: **XXX** (X.XX req/s)
- イテレーション: **XX** (X.XX iter/s)

**チェック結果**

- ✗ 2xx: **XX%** (XX/XX) - 成功したリクエスト
- ✓ 2xx or 429: **XX%** (XX/XX) - 429エラーは許容範囲として扱う

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得（各イテレーション）
2. `POST /quizzes/new` - クイズを生成（ランダムなパラメータ）

**ネットワーク**

- 受信データ: **XX kB** (XXX B/s)
- 送信データ: **XX kB** (XXX B/s)

**解釈 / ボトルネック**

- Lambda環境でのクイズ生成処理のパフォーマンス評価
- コールドスタートの影響: 初回リクエストで +XXXms の遅延
- API Gateway + Lambda + VPC + RDS のネットワークレイテンシを考慮
- クイズ生成は CPU/DB 負荷が高い処理のため、Lambda 256MB での処理性能が課題
- RDS t4g.micro での DB 接続プール（5接続）の影響
- 同時実行制限（409 another quiz is running）の影響
- クォータ制限（429エラー）の影響
- サーバーリソース不足（500エラー）の可能性
- **推奨対応**:
  - Lambda メモリサイズの増加検討（256MB → 512MB など）
  - 同時実行制限の確認と調整
  - クイズ生成処理の最適化（N+1 クエリの排除など）
  - エラーハンドリングの改善

**再現手順**

```bash
# 前提: Lambda環境がデプロイ済み
export BASE_URL="https://xxxx.execute-api.ap-northeast-1.amazonaws.com/prod"
export PROFILE="pr"  # または "nightly"
k6 run bench/k6/b/quiz_new.js
```

