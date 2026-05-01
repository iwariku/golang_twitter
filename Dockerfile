FROM golang:1.25.5

WORKDIR /app

# 依存関係
# Dockerのキャッシュを利用し、buildを早くする。そのため依存関係のファイルのみコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコピー
COPY . .

# ビルド
RUN go build -o app .

# 実行
CMD ["./app"]