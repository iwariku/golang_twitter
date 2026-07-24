FROM golang:1.25.5

WORKDIR /app

ENV GO111MODULE=on

# 依存関係
# Dockerのキャッシュを利用し、buildを早くする。そのため依存関係のファイルのみコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコピー
COPY . .

RUN go list ./...

# app と migrate の2つのバイナリを同じイメージに入れる。
# migrate用に別タスク定義を作る二重化を避け、1つのタスク定義を使い回すため。
# 通常はappを起動し、migrate時はcommandを./migrateに上書きするだけにする。
RUN go build -o app main.go
RUN go build -o migrate ./db/migrations

CMD ["./app"]