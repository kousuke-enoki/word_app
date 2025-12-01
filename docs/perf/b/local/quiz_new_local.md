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
- Executor: `ramping-vus`
- ステージ: `20s → 2VU`, `1m → 5VU`, `1m30s → 5VU`, `20s → 0VU`（PR プロファイル）
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
  - `http_req_duration{endpoint:quiz_new}: p(95)<200ms`
  - `http_req_failed{endpoint:quiz_new,status:!429}: rate<0.01`（429エラーはクォータ制限のため許容）

**結果（要約）**

- SLO: **未達成** ❌
- チェック成功率: **48.12%** (77/160)
- エラー率: **48.12%** (77/160) - quiz_new エンドポイントのみ
  - 主なエラー: 500 internal error、409 another quiz is running
- 実行時間: **約 3 分 13 秒**
- イテレーション: **80** (0.41 iter/s)
- 成功したイテレーション: **3 件**（2xx ステータスのみ）

**エンドポイント別メトリクス**

| Endpoint          | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | RPS | SLO |
| ----------------- | -------: | -------: | -------: | ---------: | --: | :-: |
| POST /quizzes/new | **921**  | **5,950**|    **-** | **48.12%** | **0.83** | ❌  |

注: このテストでは各イテレーションで複数のエンドポイントを順次呼び出します（test-login、quiz_new）

**HTTP メトリクス**

- http_req_duration (quiz_new):
  - avg: **1.47s**
  - min: **529.26ms**
  - med: **921.05ms**
  - max: **12.29s**
  - p(90): **2.44s**
  - p(95): **5.95s**
- http_req_failed (quiz_new): **48.12%** (77/160)
  - 429エラー: **0 件**（クォータ制限によるエラーなし）
  - 主なエラー: 500 internal error、409 another quiz is running
- http_reqs: **160** (0.83 req/s)
- イテレーション: **80** (0.41 iter/s)

**チェック結果**

- ✗ 2xx: **3%** (3/80) - 成功したリクエストはわずか 3 件
- ✓ 2xx or 429: **92%** (74/80) - 429エラーは許容範囲として扱う

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得（各イテレーション）
2. `POST /quizzes/new` - クイズを生成（ランダムなパラメータ）

**ネットワーク**

- 受信データ: **53 kB** (275 B/s)
- 送信データ: **49 kB** (251 B/s)

**解釈 / ボトルネック**

- SLO 未達成：p95 latency が **5.95s** で 200ms を大きく超えています
- エラー率は **48.12%** と非常に高く、システムが正常に動作していません
- 成功したリクエストは **3 件のみ**（全体の 3%）で、ほとんどがエラーです
- 主なエラーは以下の通り:
  - **500 internal error**: サーバー内部エラーが複数回発生
  - **409 another quiz is running**: 同時実行制限により、複数のクイズが同時に実行されている
- 平均レイテンシは **1.47s** と非常に高く、クイズ生成処理に問題があります
- 中央値（p50）は **921.05ms**、最大レイテンシは **12.29s** と極めて高い値です
- イテレーション実行時間の平均は **9.02s** で、ThinkTime（5-10秒）を含めても処理が重いことを示しています
- **同時実行の問題**: 409エラーが発生していることから、クイズ生成の同時実行制限が原因の可能性があります
- **サーバーリソース不足**: 500エラーの発生から、CPU/DB リソースが不足している可能性があります
- **推奨対応**:
  - 同時実行制限の確認と調整
  - サーバーリソース（CPU、メモリ、DB接続プール）の確認
  - クイズ生成処理の最適化
  - エラーハンドリングの改善

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
export BASE_URL="http://localhost:8080"
export PROFILE="pr"
k6 run bench/k6/b/quiz_new.js
```
