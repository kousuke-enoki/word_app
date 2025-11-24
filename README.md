# 初回

.env を用意して./backend, ./frontend それぞれの直下に置く

docker compose build

bash docker.sh up dev

bash docker.sh import dev

## 環境起動コマンド

docker.sh で振り分け

# 開発環境 sh コマンド

bash docker.sh up dev # LocalStack + DynamoDB + DB + Backend + Frontend を起動
bash docker.sh down dev
bash docker.sh exec backend dev
bash docker.sh exec frontend dev
bash docker.sh import dev

# LocalStack DynamoDB の初期化

# docker.sh up dev 実行時に自動的に初期化されます

# 手動で実行する場合は以下のコマンド:

LOCALSTACK_ENDPOINT=http://localhost:4566 bash scripts/init_localstack.sh

# 本番環境

bash docker.sh up prod
bash docker.sh down prod
bash docker.sh exec backend prod
bash docker.sh exec frontend prod
bash docker.sh import prod

## docker キャッシュ削除して起動

docker compose --env-file backend/.env.development down --volumes --rmi all
docker compose build --no-cache
bash docker.sh up dev

# 実行中のコンテナを確認

docker ps

# PostgreSQL コンテナに接続

docker compose exec -it db psql -U postgres -d db
docker compose --env-file backend/.env.development exec -it db psql -U postgres -d db
bash docker.sh db dev

# テーブルの一覧を表示

\dt

# テーブルの内容を表示

SELECT \* FROM users;

# カラムの内容を表示

\d users

## ent generate

# スキーマを作成

ent/schema で作成

# generate (スキーマ作ったら)

go generate ./ent

# モック作成(mockery) (v3 推奨、v2 は何故か使用できなくなった)

mockery(コンテナ内で)
go install github.com/vektra/mockery/v3@v3.4.0

go install github.com/vektra/mockery/v2@v2.43.2

# mockery v3 使用方法

.mockery.yml に、新規 interface のパッケージを追加

bash docker.sh exec backend

ルートで
mockery

# interfaces があるディレクトリで（v2）

mockery --name=UserClient --output=./../mocks

# goimport

cd backend
go install golang.org/x/tools/cmd/goimports@latest
goimports -w -local word_app/backend src/

# golangci-lint run --verbose(ライブラリ関係はスキップするように設定)

bash docker.sh exec backend
golangci-lint run --verbose

# コミット前の確認

bash docker.sh exec backend

goimports -w -local word_app/backend src/

golangci-lint run --verbose

go test ./...

## フロント

# eslint

cd frontend
npm run lint

### 自動修正

npm run lint:fix

# フロントエンドライブラリインストール

cd frontend
npm install react-i18next i18next --save

# テスト実行

cd frontend
npm test

# 品詞の詳細

1 名詞（noun）
2 代名詞（pronoun）
3 動詞（verb）
4 形容詞（adjective）
5 副詞（adverb）
6 助動詞（auxiliary verb）
7 前置詞（preposition）
8 冠詞（article）
9 間投詞（interjection）
10 接続詞（conjunction）
11 慣用句
12 その他

# 辞書 import

# Makefile 用ダウンロードコマンド

make download-dict

# もしくはこちらの URL からダウンロード

https://github.com/scriptin/jmdict-simplified/releases/download/3.6.1%2B20250421122348/jmdict-eng-3.6.1+20250421122348.json.zip

# backend/assets/ に jmdict.json をおく

# import 用コマンド （ダウンロード後にこのコマンドで db にインポート）

bash docker.sh import dev
bash docker.sh import prod

jmdict-simplified
Update checksum files @ 2025-04-07 12:33:48 UTC

https://github.com/scriptin/jmdict-simplified/tree/master?search=1

# jmdict ライセンス

CC‑BY‑SA 4.0 / WordNet License
