package main

import (
	"database/sql"
	"embed"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // sql.Open("pgx", ...)用にpgxドライバを登録する
)

// SQLをバイナリに同梱し、実行時のファイル配置ミスを構造的に防ぐ。
// go:embedは同ディレクトリ以下しか対象にできないため、コードをSQLと同じ場所に置いている。

// 20行目の //go:embed はコメントではなくコンパイラへの命令。消すと埋め込みが行われず動かない。

//go:embed *.sql
var migrationsFS embed.FS

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	migrationSource, err := iofs.New(migrationsFS, ".")
	if err != nil {
		log.Fatalf("マイグレーションの読み込みに失敗しました: %v", err)
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("DB接続に失敗しました: %v", err)
	}
	defer db.Close()

	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("マイグレーションドライバの初期化に失敗しました: %v", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", migrationSource, "postgres", migrationDriver)
	if err != nil {
		log.Fatalf("マイグレーションの初期化に失敗しました: %v", err)
	}

	// 未適用のマイグレーションをすべて適用する。
	// 適用済みで変更がなければ ErrNoChange（＝何もしない）で正常終了する。
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("マイグレーションに失敗しました: %v", err)
	}

	log.Println("migration completed")
}
