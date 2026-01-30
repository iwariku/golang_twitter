FROM golang:1.25.5
WORKDIR /app

# Dockerのキャッシュを利用し、buildを早くする。そのため依存関係のファイルのみコピー
COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/air-verse/air@latest

COPY . .
CMD ["air", "-c", ".air.toml"]