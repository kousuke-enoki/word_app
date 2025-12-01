# 目的について

## ベンチマーク方針（学習用環境としての用途）

### 想定ユーザー数

- 同時アクティブユーザー: 5-10 人
- ピーク時: 10-20 人
- 1 ユーザーあたり: 1-2 req/s

### 環境制約

- DB 接続プール: 5 接続（コスト抑制のため）
- Lambda: 256MB, 30 秒タイムアウト
- RDS: t4g.micro（最小サイズ）

### ベンチマーク条件

- 層 B (PR): 最大 5 VU, 5 req/s
- 層 B (Nightly): 最大 10 VU, 10 req/s
- 層 C: 最大 10 VU, 5-10 req/s

### 目的

- 学習用かつポートフォリオ用途として適切な負荷での動作確認
- コストを抑えつつ、基本的なパフォーマンスを保証
- 高負荷テストは実施しない（想定外の負荷のため）

## 層 A：スモーク（全 API）

目的：網羅＆回帰検知

実行：ローカル＆ステージングで各 1 リクエスト ×1〜2VU×30 秒

判定：成功率 100%（http_req_failed rate=0 近似）

CI：PR で毎回

## 層 B：パフォーマンス 4 本（SLO 厳守）

GET /words?search=（検索：読み多・索引効く）

POST /quizzes（生成：CPU/DB 負荷）

POST /auth/login（認証：外部 I/O ほぼ無し、基準線）

（任意）POST /registered_words（書き込み：ユースケース感）

実行：ステージング中心（ローカルでも確認できるようにする）

PR：最大 5 VU / 約 3 分（20s→1m→1m30s→20s のステージ、5 req/s）

Nightly：最大 10 VU / 約 4 分（20s→1m→2m→20s のステージ、10 req/s）

閾値（全シナリオ共通）：

http_req_duration: p(95)<200

http_req_failed: rate<0.01

CI：PR で短縮版＋ SLO ゲート、Nightly で本番相当

## 層 C：スポット（任意）

コールドスタート比較：アクセス断 30〜60 分後に 1 リク →p95 を記録、Provisioned Concurrency 有/無で差分

レート制限耐性：1IP/1 ユーザでスパイク →429 が予測通り出るか

DB 最適化の Before/After：EXPLAIN ＋ k6 数値を README に貼る

実行：手動 or Nightly のみ（PR では走らせない）

# 1. k6 の入れ方

## WSL/Ubuntu 内

### apt で入れる

sudo gpg -k >/dev/null || sudo apt install -y gnupg
sudo apt update
sudo apt install -y ca-certificates curl
curl -s https://dl.k6.io/key.gpg | sudo gpg --dearmor -o /usr/share/keyrings/k6-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt update && sudo apt install -y k6
k6 version

## Windows(ネイティブ)

choco install k6 -y # Chocolatey

### または Scoop: scoop install k6

k6 version

## macOS

brew install k6
k6 version

# A local 全 API の動作確認

BASE_URL=http://localhost:8080 SMOKE_ENV=local k6 run bench/k6/a/smoke_all.js

# B 各 API 詳細向けの測定

## ローカル環境向け

ローカル環境では既存の処理が使用されます（閾値: p95<200ms）。

BASE_URL=http://localhost:8080 PROFILE=pr TEST_EMAIL=demo@example.com TEST_PASSWORD=Secret-k6 k6 run bench/k6/b/sign_in.js

BASE_URL=http://localhost:8080 PROFILE=pr SEARCH_Q=able SEARCH_SORT=name k6 run bench/k6/b/register_word.js

BASE_URL=http://localhost:8080 PROFILE=pr k6 run bench/k6/b/quiz_new.js

BASE_URL=http://localhost:8080 PROFILE=pr SEARCH_Q=test SEARCH_SORT=name k6 run bench/k6/b/search_words.js

## Lambda 環境向け

Lambda 環境では自動検出され、以下の最適化が適用されます：

- ウォームアップフェーズ追加（コールドスタート対策）
- 閾値緩和（p95<1000ms、コールドスタートとネットワークレイテンシを考慮）
- setup()でのリトライ処理（タイムアウト対策）

環境検出は `BASE_URL` に `amazonaws.com` が含まれている場合に自動で行われます。
明示的に指定する場合は `IS_LAMBDA=true` または `IS_LAMBDA=false` を設定してください。

BASE_URL="https://xxxx.execute-api.ap-northeast-1.amazonaws.com/prod" PROFILE=pr \
TEST_EMAIL=demo@example.com TEST_PASSWORD='K6passw0rd!' \
k6 run bench/k6/b/sign_in.js

BASE_URL="https://xxxx.execute-api.ap-northeast-1.amazonaws.com/prod" PROFILE=pr SEARCH_Q=test SEARCH_SORT=name \
k6 run bench/k6/b/search_words.js

BASE_URL="https://xxxx.execute-api.ap-northeast-1.amazonaws.com/prod" PROFILE=pr \
k6 run bench/k6/b/quiz_new.js

BASE_URL="https://xxxx.execute-api.ap-northeast-1.amazonaws.com/prod" PROFILE=pr SEARCH_Q=able SEARCH_SORT=name \
k6 run bench/k6/b/register_word.js

# C

# 1 回だけフォルダ作成

mkdir -p bench/k6/out

# cold/warm

BASE_URL=http://localhost:8080 \
npm --prefix bench run k6:c:cold

# rate-limit

<!-- BASE_URL=http://localhost:8080 \
npm --prefix bench run k6:c:rate -->

# DB before/after（ラベルは --summary-export ファイル名で区別する運用でも OK）

BASE_URL=http://localhost:8080 LABEL=before_idx \
npm --prefix bench run k6:c:db:before

BASE_URL=http://localhost:8080 LABEL=after_idx \
npm --prefix bench run k6:c:db:after

# RPS–p95 の時系列グラフを自動生成

## k6 からサンプル時系列を吐く（--out json）

BASE_URL=$STG_URL PROFILE=pr \
k6 run --out json=bench/k6/out/search_timeseries.json bench/k6/b/search_words.js

BASE_URL=http://localhost:8080 PROFILE=pr \
k6 run --out json=bench/k6/out/search_timeseries.json bench/k6/b/search_words.js

node bench/tools/mkcharts.js bench/k6/out/search_timeseries.json bench/k6/out/search_rps_p95.html

→ できた HTML を README から相対リンク、あるいはスクショ PNG を一緒にコミットなどする
