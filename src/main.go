package main

func main() {
	load_config()
	db := db_init()
	bot_start(db)
}
