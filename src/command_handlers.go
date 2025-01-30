package main

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartCommand(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {

	user := User{}

	if data == "return" {
		user.Id = update.CallbackQuery.From.ID
		user.State = "menu"
	} else {
		user.Id = update.Message.From.ID
		user.State = "menu"
	}

	err := user.Insert(db)

	if err != nil {
		user.Update(db)
	}

	msg_Text := "Main menu"
	msg_ReplyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Create habit ⚒️", "create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("My habits 📋", "myhabits"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Statistic 📈", "statistic"),
		),
	)

	if data != "return" {
		msg := tgbotapi.NewMessage(user.Id, "")
		msg.Text = msg_Text
		msg.ReplyMarkup = msg_ReplyMarkup

		_, err = bot.Send(msg)

		if err != nil {
			log.Panicln(err)
		}
	} else {
		msg_edit := tgbotapi.NewEditMessageTextAndMarkup(user.Id,
			update.CallbackQuery.Message.MessageID,
			msg_Text,
			msg_ReplyMarkup,
		)

		_, err = bot.Send(msg_edit)

		if err != nil {
			log.Panicln(err)
		}
	}
}

func CreateCommandRouter() *Router {
	router := NewRouter()

	router.register("start", StartCommand)
	router.register("menu", StartCommand)

	return router
}
