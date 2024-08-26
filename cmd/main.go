package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ShaynaSegal45/phonebook-api/contactsmanaging"
	mysql "github.com/ShaynaSegal45/phonebook-api/sql"
)

func main() {
	db := initializeDatabase()
	defer db.Close()

	repo := mysql.NewContactsRepo(db)
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
