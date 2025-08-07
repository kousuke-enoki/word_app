# backend/dockerfiles/Lambda.Dockerfile

# --- build stage -------------------------------------------------------------
FROM --platform=linux/amd64 golang:1.24.3 AS build
WORKDIR /src

# 依存キャッシュ
COPY go.mod go.sum ./
RUN go mod download

# アプリ全体コピー
COPY . .

# 必要なら ent のコード生成（generate.go の go:generate を利用）
# 生成済みコードを Git 管理しているなら以下は削ってもOK
RUN go generate ./ent

# Lambda Custom Runtime は /bootstrap 実行が規約
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /bootstrap ./cmd/server

# --- runtime stage (AWS Lambda provided.al2) ---------------------------------
FROM --platform=linux/amd64 public.ecr.aws/lambda/provided:al2

# 実行ファイルのみコピー
COPY --from=build /bootstrap /bootstrap

# デフォルトは /bootstrap 実行
# 環境変数は Dockerfileに書かず、Lambda関数の設定（CDK）で注入します
# CMD は不要（/bootstrap が既定）
