# ビルドステージ
FROM golang:1.24-alpine AS builder

# ビルド時引数を受け取る
ARG VERSION=dev
ARG GIT_COMMIT=unknown

# 作業ディレクトリ設定
WORKDIR /app

# 依存関係ファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド（ldflags で VERSION と GIT_COMMIT を埋め込み）
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags "-w -s -X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT}" \
    -o main .

# 実行ステージ
FROM alpine:latest

# セキュリティ：非rootユーザー作成
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# CA証明書インストール（HTTPSリクエスト用）
RUN apk --no-cache add ca-certificates

WORKDIR /app

# ビルドステージからバイナリをコピー
COPY --from=builder /app/main .

# 非rootユーザーに切り替え
USER appuser

# ポート公開
EXPOSE 8080

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# 実行
CMD ["./main"]