# backend/dockerfiles/Import.Dockerfile

# --- build stage -------------------------------------------------------------
FROM --platform=linux/amd64 golang:1.25.1 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 必要なら ent のコード生成
RUN go generate ./ent

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /import_dict ./cmd/import_dict

# --- runtime stage (distroless) ---------------------------------------------
FROM --platform=linux/amd64 gcr.io/distroless/static-debian12
WORKDIR /app

COPY --from=build /import_dict /usr/local/bin/import_dict

# デフォルトエントリポイント
ENTRYPOINT ["/usr/local/bin/import_dict"]
# 例: 引数がなければ usage を表示、またはアプリ側でデフォルトを定義
# CMD ["-file=/data/jmdict.json", "-workers=4"]
