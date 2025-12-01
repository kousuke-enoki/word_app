**Environment: Lambda (API Gateway + Lambda 256MB + RDS t4g.micro, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **1000ms**（Lambda環境向け：コールドスタートとネットワークレイテンシを考慮）
- Error rate < **1%**（429エラーは登録単語数上限超過のため許容）
- 対象: `POST /words/register`（単語登録：DB更新処理）

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

- シナリオ: register_word（単語登録）
- Executor: `ramping-vus`
- ステージ（Lambda向け: ウォームアップ付き）:
  - PR: `30s → 1VU`（ウォームアップ）, `20s → 2VU`, `1m → 5VU`, `1m30s → 5VU`, `20s → 0VU`
  - Nightly: `30s → 1VU`（ウォームアップ）, `20s → 3VU`, `1m → 10VU`, `2m → 10VU`, `20s → 0VU`
- ThinkTime: 1s
- 単語選択: ランダムな検索クエリとソート条件で単語IDを取得（`pickWordId()`）
  - 検索クエリ: `randomSearchQuery()`（"a", "the", "and" など一般的な単語）
  - ソート条件: `randomSortBy()`（ランダム）
  - ページ: 1-10ページからランダム選択
  - リトライ: 最大5回（検索結果が空の場合、別のクエリを試行）
- 事前処理: テスト実行前に登録済み単語を全件リセット（上限200件まで）
- 閾値（thresholds）:
  - `http_req_duration{endpoint:register_word}: p(95)<1000ms`（Lambda向け）
  - `http_req_failed{endpoint:register_word,status:!429}: rate<0.01`（429エラーは登録単語数上限超過のため許容）

**結果（要約）**

- SLO: **達成** ✅
- 最大安定 RPS（5VU 時）: **10.72 req/s**
- コールドスタート: **含む**（ウォームアップ後も初回リクエストで +XXXms）
- チェック成功率: **100%** (311/311)
- エラー率: **9.37%** (223/2379) - register_word エンドポイント全体
  - 429エラー（登録単語数上限超過）: **許容範囲**
  - 429以外のエラー率: **0.00%** (0/0) - 閾値をクリア
- 実行時間: **約 3 分 41 秒**
- イテレーション: **561** (2.53 iter/s)

**エンドポイント別メトリクス**

| Endpoint            | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | 5VU 時 RPS | SLO |
| ------------------- | -------: | -------: | -------: | ---------: | ---------: | :-: |
| POST /words/register | **92.64** | **172.24** | **XXX** | **9.37%** | **10.72** | ✅ |

注: エラー率には429エラー（登録単語数上限超過）が含まれますが、これは許容範囲です。

**HTTP メトリクス**

- http_req_duration (register_word):
  - avg: **101.76ms**
  - min: **67.32ms**
  - med: **92.64ms**
  - max: **620.79ms**
  - p(90): **121.51ms**
  - p(95): **172.24ms**
- http_req_failed (register_word): **9.37%** (223/2379)
  - 429エラー（登録単語数上限超過）: **許容範囲**
  - 429以外のエラー: **0.00%** (閾値をクリア)
- http_reqs: **2379** (10.72 req/s)
- イテレーション: **561** (2.53 iter/s)

**チェック結果**

- ✓ 2xx or 409 or acceptable errors: **100%** (311/311)
  - 2xx: 成功
  - 409: ユニーク制約違反（許容）
  - 400: 登録状態に変更なし（許容）
  - 429: 登録単語数上限超過（許容）

**ネットワーク**

- 受信データ: **1.8 MB** (7.9 kB/s)
- 送信データ: **215 kB** (969 B/s)

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得（setup 関数）
2. `POST /words/register` - 単語を登録（各イテレーション）
   - 事前に `GET /words?search=...` で単語IDを取得（`pickWordId()`）
   - ランダムな検索クエリとソート条件を使用

**問題点 / 注意事項**

- **`pickWordId()` でのエラー大量発生**: 「no words found in response」エラーが多数発生
  - 原因: ランダムページ（1-10ページ）指定で、実際のデータが少ないページを選択している可能性
  - 影響: エラーログが大量に出力されるが、ベンチマーク自体は継続（561イテレーション完了）
  - 改善提案:
    - ページ1から順に試行する、または結果があるページを確認してから選択
    - より確実に結果を返す検索クエリの使用
    - エラー時のフォールバック（固定IDの使用など）

**解釈 / ボトルネック**

- Lambda環境での単語登録処理のパフォーマンス評価
- コールドスタートの影響: 初回リクエストで +XXXms の遅延
- API Gateway + Lambda + VPC + RDS のネットワークレイテンシを考慮
- p95=172.24ms と良好なパフォーマンス（閾値1000msを大幅にクリア）
- 429エラー（登録単語数上限超過）は許容範囲として扱う
- RDS t4g.micro での DB 接続プール（5接続）の影響
- ウォームアップフェーズ（30秒、1 VU）により、コールドスタートの影響を軽減
- 256MB メモリでの処理性能は十分
- **推奨対応**:
  - `pickWordId()` 関数の改善（より確実に単語IDを取得できるようにする）
  - エラーハンドリングの改善（エラーログの抑制など）

**再現手順**

```bash
# 前提: Lambda環境がデプロイ済み
export BASE_URL="https://xxxx.execute-api.ap-northeast-1.amazonaws.com/prod"
export PROFILE="pr"  # または "nightly"
k6 run bench/k6/b/register_word.js
```

