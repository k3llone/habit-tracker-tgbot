package main

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartTextHandler(text string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	msg := tgbotapi.NewMessage(update.Message.From.ID, "")
	msg.Text = "Menu state"

	_, err := bot.Send(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func CreateTextRouter() *TextRouter {
	router := NewTextRouter()

	router.register("menu", StartTextHandler)

	return router
}
