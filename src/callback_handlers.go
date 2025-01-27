package main

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func OneData(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Data)
	msg.Text = "You press button 1"

	_, err := bot.Send(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func CreateCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	user := User{}
	user.Load(update.CallbackQuery.From.ID, db)

	user.State = "create_name"

	user.Update(db)

	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "")
	msg.Text = "Enter name (no more 20 characters)"

	_, err := bot.Send(msg)

	if err != nil {
		log.Println(err)
	}
}

func MtHabitsCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {

}

func StatisticCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {

}

func CreateCallbackRouter() *Router {
	router := NewRouter()

	router.register("create", CreateCallbackHandler)
	router.register("habits", MtHabitsCallbackHandler)
	router.register("statistic", StatisticCallbackHandler)

	return router
}
