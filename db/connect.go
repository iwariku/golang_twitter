package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDB(ctx context.Context) (*pgx.Conn, *Queries) {
	databaseUrl := os.Getenv("DATABASE_URL")

	var connStr string

	if databaseUrl == "" {
		// ローカル（Docker）
		user := os.Getenv("POSTGRES_USER")
		pass := os.Getenv("POSTGRES_PASSWORD")
		host := "postgres"
		port := "5432"
		dbname := os.Getenv("POSTGRES_DB")

		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, dbname)
	} else {
		// Render（本番）
		connStr = databaseUrl
	}

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("DBに接続できませんでした: %v", err)
	}

	fmt.Println("DB接続に成功しました")

	return conn, New(conn)
}
