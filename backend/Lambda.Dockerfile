FROM --platform=linux/amd64 golang:1.24.3 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /main ./cmd/health

# Go ランタイムベース
FROM --platform=linux/amd64 public.ecr.aws/lambda/go:1
# main は /var/task/main になる
COPY --from=build /main /var/task/main
# こちらは CMD で "main" を指定
CMD ["main"]
