# backend/dockerfiles/Builder.Dockerfile
FROM golang:1.25.1 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# ── ★ マイグレーション・コード生成 etc…
RUN go generate ./ent

# ① API
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server
# ② dict-import
RUN CGO_ENABLED=0 go build -o /out/import_dict ./cmd/import_dict
