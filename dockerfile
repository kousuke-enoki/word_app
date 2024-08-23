# ベースイメージを指定
FROM golang:latest AS word_app

# コンテナ内の作業ディレクトリを設定
WORKDIR /app

# ローカルのソースコードをコンテナにコピー
COPY ./backend ./

# 必要なパッケージをインストール
RUN go mod download

# アプリケーションをビルド
RUN go build -o main .

# フロントエンドのビルドステージ
FROM node:latest AS frontend

# コンテナ内の作業ディレクトリを設定
WORKDIR /app/frontend

# ローカルのソースコードをコンテナにコピー
COPY ./frontend ./

# 必要なパッケージをインストール
RUN npm install

# フロントエンドアプリケーションをビルド
RUN npm run build

# 最終ステージ
FROM golang:latest

# バックエンドのビルド成果物をコピー
COPY --from=backend /app/main /app/main

# フロントエンドのビルド成果物をコピー
COPY --from=frontend /app/frontend/build /app/frontend/build

# コンテナ内の作業ディレクトリを設定
WORKDIR /app

# 公開予定のコンテナのポートを明示
EXPOSE 8080

# アプリケーションを実行
CMD ["./main"]