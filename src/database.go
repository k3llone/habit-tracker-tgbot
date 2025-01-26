package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func db_init() *sql.DB {
	db, err := sql.Open("sqlite3", "data.db")

	if err != nil {
		log.Panicln(err)
	}

	return db
}
