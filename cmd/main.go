package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ShaynaSegal45/phonebook-api/contactsmanaging"
	sqldb "github.com/ShaynaSegal45/phonebook-api/sql"
)

func main() {
	db := initializeDatabase()
	rdb := initRedisClient()
	defer db.Close()
	defer rdb.Close()

	repo := sqldb.NewContactsRepo(db, rdb)
	service := contactsmanaging.NewService(repo)
	router := contactsmanaging.NewHTTPHandler(service)

	startServer(router)
}

func initializeDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./contacts.db")
	if err != nil {
		log.Fatalf("could not connect to database: %v\n", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS contacts (
		id TEXT PRIMARY KEY,
		firstname TEXT,
		lastname TEXT,
		address TEXT,
		phone TEXT
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("could not create contacts table: %v\n", err)
	}

	return db
}

func startServer(router http.Handler) {
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}

func initRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}
	fmt.Println(pong)

	return rdb
}
