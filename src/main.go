package main

import "log"

func main() {
	load_config()

	db, err := db_init()

	if err != nil {
		log.Panicln(err)
	}

	db_create(db)

	bot_start(db)
}
