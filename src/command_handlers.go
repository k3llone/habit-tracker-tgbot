package main

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartCommand(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {

	user := User{}
	user.Id = update.Message.From.ID
	user.State = "menu"

	err := user.Insert(db)

	if err != nil {
		user.Update(db)
	}

	msg := tgbotapi.NewMessage(update.Message.From.ID, "")
	msg.Text = "Main menu"

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Create habit", "create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("My habits", "habits"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Statistic", "statistic"),
		),
	)

	_, err = bot.Send(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func CreateCommandRouter() *Router {
	router := NewRouter()

	router.register("start", StartCommand)

	return router
}
