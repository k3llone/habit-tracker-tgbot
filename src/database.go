package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func db_init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "data.db")

	if err != nil {
		return nil, err
	}

	return db, nil
}

func db_create(db *sql.DB) {
	db.Exec("CREATE TABLE IF NOT EXISTS \"HabitComplete\" (\"id\" INTEGER UNIQUE, \"habit\" INTEGER, \"date\" TEXT, PRIMARY KEY(\"id\" AUTOINCREMENT));")
	db.Exec("CREATE TABLE IF NOT EXISTS \"HabitMenus\" (\"id\" INTEGER UNIQUE, \"user\" INTEGER, \"pages\" INTEGER, \"cpage\" INTEGER, \"habits\" TEXT, PRIMARY KEY(\"id\" AUTOINCREMENT));")
	db.Exec("CREATE TABLE IF NOT EXISTS \"habits\" ( \"id\"	integer unique, \"user\" integer, \"name\" text, \"remtime\" text, primary key(\"id\" autoincrement));")
	db.Exec("CREATE TABLE IF NOT EXISTS \"users\" ( \"id\" integer unique, \"states\" text, \"createhabit\" integer);")
}
