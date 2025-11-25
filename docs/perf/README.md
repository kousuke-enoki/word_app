# Performance Test Results

## 層分けの目的

パフォーマンステストは、目的に応じて 3 つの層に分けられています：

- **A 層: スモークテスト** - 回帰検知、最速 30 秒で全 API の動作確認
- **B 層: SLO 検証** - 代表 API× 負荷プロファイルで p95/error 検証
- **C 層: スポットテスト** - 仮説検証（cold/warm、DB 最適化など）

## SLO & Thresholds

**SLO（目標）**:

- p95 latency < **200ms**
- Error rate < **1%**

**Thresholds（k6 閾値）**:

- B 層・C 層: `http_req_duration: p(95)<200ms`, `http_req_failed: rate<0.01`
- A 層: `http_req_failed: rate<0.01` のみ（p95 閾値なし）
  - 理由: スモークテストは動作確認が目的のため、レイテンシ閾値は不要

## 命名規則

テスト結果ファイルは環境プレフィックスで区別します：

- `*_local.md` - ローカル環境（Docker、ローカル DB）
- `*_demo.md` - デモ環境（Lambda/本番相当、将来追加）
- `*_staging.md` - ステージング環境（将来追加）
- `*_prod.md` - 本番環境（将来追加）

## テスト結果一覧

### A 層: スモークテスト

- [全 API 動作確認 (smoke_local)](a/smoke_local.md) - 全エンドポイントの動作確認

### B 層: パフォーマンステスト

- [認証 (sign_in_local)](b/sign_in_local.md) - `POST /users/sign_in`
- [単語検索 (search_words_local)](b/search_words_local.md) - `GET /words?q=`
- [クイズ生成 (quiz_new_local)](b/quiz_new_local.md) - `POST /quizzes/new`
- [単語登録 (register_word_local)](b/register_word_local.md) - `POST /words/register`

### C 層: スポットテスト

- [コールドスタート比較 (cold_warm_local)](c/cold_warm_local.md) - コールドスタート vs ウォームスタート
- [DB 最適化前 (db_before_local)](c/db_before_local.md) - インデックス最適化前の性能
- [DB 最適化後 (db_after_local)](c/db_after_local.md) - インデックス最適化後の性能

## Before/After 比較

DB 最適化の効果を可視化：

![DB Index Before/After](db_index_before_after.png)

詳細は [db_before_local.md](c/db_before_local.md) と [db_after_local.md](c/db_after_local.md) を参照してください。
