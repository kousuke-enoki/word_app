# 目的について

## 層 A：スモーク（全 API）

目的：網羅＆回帰検知

実行：ローカル＆ステージングで各 1 リクエスト ×1〜2VU×30 秒

判定：成功率 100%（http_req_failed rate=0 近似）

CI：PR で毎回

## 層 B：パフォーマンス 4 本（SLO 厳守）

GET /words?q=（検索：読み多・索引効く）

POST /quizzes（生成：CPU/DB 負荷）

POST /auth/login（認証：外部 I/O ほぼ無し、基準線）

（任意）POST /registered_words（書き込み：ユースケース感）

実行：ステージング中心（ローカルでも確認できるようにする）

PR：VU 10 / 4 分（1→10→0 のステージ）

Nightly：VU 30〜50 / 6 分

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

BASE_URL=http://localhost:8080 PROFILE=pr TEST_EMAIL=demo@example.com TEST_PASSWORD=Secret-k6 k6 run bench/k6/b/sign_in.js

BASE_URL=http://localhost:8080 PROFILE=pr SEARCH_Q=able SEARCH_SORT=name k6 run bench/k6/b/register_word.js

BASE_URL=http://localhost:8080 PROFILE=pr k6 run bench/k6/b/quiz_new.js

BASE_URL=http://localhost:8080 PROFILE=pr SEARCH_Q=test SEARCH_SORT=name k6 run bench/k6/b/search_words.js

# B lambda 向け

BASE_URL="https://xxxx.execute-api.ap-northeast-1.amazonaws.com/prod" PROFILE=pr \
TEST_EMAIL=demo@example.com TEST_PASSWORD='K6passw0rd!' \
k6 run bench/k6/b/sign_in.js

BASE_URL="https://.../prod" PROFILE=pr SEARCH_Q=test SEARCH_SORT=name \
k6 run bench/k6/b/search_words.js

BASE_URL="https://.../prod" PROFILE=pr \
k6 run bench/k6/b/quiz_new.js

BASE_URL="https://.../prod" PROFILE=pr SEARCH_Q=able SEARCH_SORT=name \
k6 run bench/k6/b/register_word.js

# C

# 1 回だけフォルダ作成

mkdir -p bench/k6/out

# cold/warm

BASE_URL=http://localhost:8080 \
npm --prefix bench run k6:c:cold

# rate-limit

BASE_URL=http://localhost:8080 \
npm --prefix bench run k6:c:rate

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
