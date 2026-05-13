package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Email    string
	Password string
	UserName string
	NickName string
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	ctx := context.Background()

	var users []User

	for i := 1; i <= 15; i++ {
		users = append(users, User{
			Email:    fmt.Sprintf("test%02d@example.com", i),
			Password: "Pass1234!",
			UserName: fmt.Sprintf("test%02d", i),
			NickName: fmt.Sprintf("Test%02d", i),
		})
	}

	// =========================================
	// users
	// =========================================

	for i := range users {

		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(users[i].Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			log.Fatal(err)
		}

		err = db.QueryRowContext(
			ctx,
			`
			INSERT INTO users (
				email,
				password,
				user_name,
				nick_name,
				is_active,
				created_at,
				updated_at
			)
			VALUES ($1,$2,$3,$4,true,NOW(),NOW())
			RETURNING id;
			`,
			users[i].Email,
			string(hashedPassword),
			users[i].UserName,
			users[i].NickName,
		).Scan(&users[i].ID)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("users created")

	// =========================================
	// tweets
	// =========================================

	tweets := []struct {
		UserID  int
		Content string
	}{
		{users[0].ID, "Go + sqlc 構成かなり好き"},
		{users[2].ID, "Render deploy 完了"},
		{users[4].ID, "Gin framework 軽くて良い"},
		{users[6].ID, "最近 PostgreSQL 勉強中"},
		{users[8].ID, "ページネーション実装できた"},
		{users[10].ID, "DM機能作成中"},
		{users[12].ID, "Tailwind CSS v4 良い感じ"},
		{users[14].ID, "Go backend 楽しい"},
	}

	var tweetIDs []int

	for _, tweet := range tweets {

		var tweetID int

		err := db.QueryRowContext(
			ctx,
			`
			INSERT INTO tweets (
				user_id,
				content,
				created_at
			)
			VALUES ($1,$2,NOW())
			RETURNING id;
			`,
			tweet.UserID,
			tweet.Content,
		).Scan(&tweetID)

		if err != nil {
			log.Fatal(err)
		}

		tweetIDs = append(tweetIDs, tweetID)
	}

	fmt.Println("tweets created")

	// =========================================
	// follows
	// =========================================

	test01 := users[0]

	for i := 1; i < len(users); i++ {

		// test01 follows everyone
		_, err := db.ExecContext(
			ctx,
			`
			INSERT INTO follows (
				follower_id,
				following_id,
				created_at,
				updated_at
			)
			VALUES ($1,$2,NOW(),NOW());
			`,
			test01.ID,
			users[i].ID,
		)

		if err != nil {
			log.Fatal(err)
		}

		// everyone follows test01
		_, err = db.ExecContext(
			ctx,
			`
			INSERT INTO follows (
				follower_id,
				following_id,
				created_at,
				updated_at
			)
			VALUES ($1,$2,NOW(),NOW());
			`,
			users[i].ID,
			test01.ID,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("follows created")

	// =========================================
	// likes
	// =========================================

	for i, tweetID := range tweetIDs {

		user := users[(i+1)%len(users)]

		_, err := db.ExecContext(
			ctx,
			`
			INSERT INTO likes (
				user_id,
				tweet_id,
				created_at
			)
			VALUES ($1,$2,NOW());
			`,
			user.ID,
			tweetID,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("likes created")

	// =========================================
	// retweets
	// =========================================

	for i, tweetID := range tweetIDs {

		user := users[(i+2)%len(users)]

		_, err := db.ExecContext(
			ctx,
			`
			INSERT INTO retweets (
				user_id,
				tweet_id,
				created_at
			)
			VALUES ($1,$2,NOW());
			`,
			user.ID,
			tweetID,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("retweets created")

	// =========================================
	// bookmarks
	// =========================================

	for i, tweetID := range tweetIDs {

		user := users[(i+3)%len(users)]

		_, err := db.ExecContext(
			ctx,
			`
			INSERT INTO bookmarks (
				user_id,
				tweet_id,
				created_at
			)
			VALUES ($1,$2,NOW());
			`,
			user.ID,
			tweetID,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("bookmarks created")

	// =========================================
	// dm group
	// =========================================

	var groupID int

	err = db.QueryRowContext(
		ctx,
		`
		INSERT INTO dm_groups (
			name,
			created_at
		)
		VALUES ($1,NOW())
		RETURNING id;
		`,
		"demo-group",
	).Scan(&groupID)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("dm group created")

	// =========================================
	// dm members
	// =========================================

	for _, user := range users[:3] {

		_, err := db.ExecContext(
			ctx,
			`
			INSERT INTO dm_group_members (
				user_id,
				dm_group_id,
				created_at
			)
			VALUES ($1,$2,NOW());
			`,
			user.ID,
			groupID,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("dm members created")

	// =========================================
	// dm messages
	// =========================================

	messages := []struct {
		UserID  int
		Message string
	}{
		{users[0].ID, "こんにちは！"},
		{users[1].ID, "DMテストです"},
		{users[2].ID, "Render 本番環境です"},
		{users[0].ID, "ページネーション確認中"},
		{users[1].ID, "Go + sqlc 構成良い感じ"},
	}

	for _, message := range messages {

		_, err := db.ExecContext(
			ctx,
			`
			INSERT INTO dm_messages (
				user_id,
				dm_group_id,
				message,
				created_at
			)
			VALUES ($1,$2,$3,NOW());
			`,
			message.UserID,
			groupID,
			message.Message,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("dm messages created")

	fmt.Println("seed completed")
}
