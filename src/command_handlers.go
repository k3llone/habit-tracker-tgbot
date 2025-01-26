package main

import (
	"database/sql"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {

	user := User{}
	user.Id = update.Message.From.ID
	user.State = make([]string, 0)
	user.State = append(user.State, "menu")

	user.Insert(db)

	msg := tgbotapi.NewMessage(update.Message.From.ID, "")
	msg.Text = fmt.Sprintf("Hello, %v!", update.Message.From.UserName)

	_, err := bot.Send(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func CreateCommandRouter() *Router {
	router := NewRouter()

	router.register("start", StartCommand)

	return router
}
