# eng_backend
docker compose build

bash start.sh development

docker compose exec backend bash


## 環境起動コマンド
 docker.shで振り分け

# 開発環境 sh コマンド
bash docker.sh up dev
bash docker.sh down dev
bash docker.sh exec backend dev
bash docker.sh exec frontend dev

# 本番環境
bash docker.sh up prod
bash docker.sh down prod
bash docker.sh exec backend prod
bash docker.sh exec frontend prod


## dockerキャッシュ削除
docker compose --env-file backend/.env.development down --volumes --rmi all
docker compose build --no-cache
bash docker.sh up dev
<!-- bash docker.sh up dev -->

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
npm run eslint

# フロントエンドライブラリインストール
cd frontend
npm install react-i18next i18next --save

# テスト実行
cd frontend
npm test


# 品詞の詳細
0 名詞（noun）
1 代名詞（pronoun）
2 動詞（verb）
3 形容詞（adjective）
4 副詞（adverb）
5 助動詞（auxiliary verb）
6 前置詞（preposition）
7 冠詞（article）
8 間投詞（interjection）
9 接続詞（conjunction）
