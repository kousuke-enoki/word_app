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