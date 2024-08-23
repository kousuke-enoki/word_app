# ベースイメージを指定
FROM golang:latest

# コンテナ内の作業ディレクトリを設定
WORKDIR /app

# ローカルのソースコードをコンテナにコピー
COPY . .

# 必要なパッケージをインストール
RUN go mod download

# アプリケーションをビルド
RUN go build -o main .

# 公開予定のコンテナのポートを明示
EXPOSE 8080

# アプリケーションを実行
CMD ["go", "run", "main.go"]
