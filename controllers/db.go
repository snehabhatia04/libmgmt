package controllers

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB struct to hold the database connection
type DB struct {
	*sqlx.DB
}

// ConnectDB initializes and returns a database connection
func ConnectDB() (*DB, error) {
	// Connect to the database
	db, err := sqlx.Connect("postgres", "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
		return nil, err
	}
	return &DB{db}, nil
}
