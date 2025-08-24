# FROM --platform=linux/amd64 golang:1.24.3 AS build
# WORKDIR /src
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
#     go build -trimpath -ldflags="-s -w" -o /main ./cmd/server

# # Go ランタイムベース
# FROM --platform=linux/amd64 public.ecr.aws/lambda/go:1
# # main は /var/task/main になる
# COPY --from=build /main /var/task/main
# # こちらは CMD で "main" を指定
# CMD ["main"]
# ---- build stage ----
    FROM --platform=linux/amd64 golang:1.24.3 AS build
    WORKDIR /src
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    
    # ← ここが肝。デフォルト server、必要なら health に切替可
    ARG TARGET=server
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -trimpath -ldflags="-s -w" -o /main ./cmd/${TARGET}
    
    # ---- runtime stage (Go Lambda 公式ランタイム) ----
    FROM --platform=linux/amd64 public.ecr.aws/lambda/go:1
    COPY --from=build /main /var/task/main
    CMD ["main"]   # ← これが“ハンドラー名”。これを使う
    
# 使い分け↓
# 本番：--build-arg TARGET=server
# 疑似ヘルス：--build-arg TARGET=health