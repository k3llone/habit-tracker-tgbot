package main

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func OneData(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Data)
	msg.Text = "You press button 1"

	_, err := bot.Send(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func CreateCallbackRouter() *Router {
	router := NewRouter()

	router.register("1", OneData)

	return router
}
