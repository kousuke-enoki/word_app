# eng_backend

docker compose build

docker compose up

docker compose exec backend bash



〇db接続方法

# 実行中のコンテナを確認
docker ps

# PostgreSQLコンテナに接続
docker compose exec -it db psql -U postgres -d db

# テーブルの一覧を表示
\dt

# テーブルの内容を表示
SELECT * FROM users;

# カラムの内容を表示
\d users




〇ent generate

# スキーマを作成
ent/schema で作成

#  generate
go generate ./ent


〇dockerキャッシュ削除
docker compose down --volumes --rmi all
docker compose build --no-cache
docker compose up


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
