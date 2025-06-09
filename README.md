# 初回
docker compose build

bash docker.sh up dev

bash docker.sh import dev

## 環境起動コマンド
 docker.shで振り分け

# 開発環境 sh コマンド
bash docker.sh up dev
bash docker.sh down dev
bash docker.sh exec backend dev
bash docker.sh exec frontend dev
bash docker.sh import dev


# 本番環境
bash docker.sh up prod
bash docker.sh down prod
bash docker.sh exec backend prod
bash docker.sh exec frontend prod
bash docker.sh import prod


## dockerキャッシュ削除して起動
docker compose --env-file backend/.env.development down --volumes --rmi all
docker compose build --no-cache
bash docker.sh up dev

# 実行中のコンテナを確認
docker ps

# PostgreSQLコンテナに接続
docker compose exec -it db psql -U postgres -d db
docker compose --env-file backend/.env.development exec -it db psql -U postgres -d db
bash docker.sh db dev

# テーブルの一覧を表示
\dt

# テーブルの内容を表示
SELECT * FROM users;

# カラムの内容を表示
\d users


## ent generate

# スキーマを作成
ent/schema で作成

#  generate (スキーマ作ったら)
go generate ./ent

# モック作成(mockery)
mockery(コンテナ内で)
go install github.com/vektra/mockery/v2@v2.43.2

interfacesがあるディレクトリで
mockery --name=UserClient --output=./../mocks

# goimport
cd backend
goimports -w -local word_app/backend src/

# golangci-lint run --verbose(ライブラリ関係はスキップするように設定)
cd backend
golangci-lint run --verbose


## フロント

#  eslint
cd frontend
npm run eslint

# フロントエンドライブラリインストール
cd frontend
npm install react-i18next i18next --save

# テスト実行
cd frontend
npm test


# 品詞の詳細
1  名詞（noun）
2  代名詞（pronoun）
3  動詞（verb）
4  形容詞（adjective）
5  副詞（adverb）
6  助動詞（auxiliary verb）
7  前置詞（preposition）
8  冠詞（article）
9  間投詞（interjection）
10 接続詞（conjunction）
11 慣用句
12 その他


# 辞書import
# Makefile用ダウンロードコマンド
make download-dict

# もしくはこちらのURLからダウンロード
https://github.com/scriptin/jmdict-simplified/releases/download/3.6.1%2B20250421122348/jmdict-eng-3.6.1+20250421122348.json.zip

# backend/assets/ にjmdict.jsonをおく

# import用コマンド （ダウンロード後にこのコマンドでdbにインポート）
bash docker.sh import dev
bash docker.sh import prod


jmdict-simplified
Update checksum files @ 2025-04-07 12:33:48 UTC

https://github.com/scriptin/jmdict-simplified/tree/master?search=1



# jmdict ライセンス
CC‑BY‑SA 4.0 / WordNet License
