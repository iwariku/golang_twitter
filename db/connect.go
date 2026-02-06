package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDB(ctx context.Context) (*pgx.Conn, *Queries) {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := "postgres"
	port := "5432"
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, dbname)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("DBに接続できませんでした: %v", err)
	}

	fmt.Println("DB接続に成功しました")
	return conn, New(conn)
}
