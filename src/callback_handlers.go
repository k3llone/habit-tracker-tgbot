package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func OneData(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Data)
	msg.Text = "You press button 1"

	_, err := bot.Send(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func TwoData(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Data)
	msg.Text = "You press button 2"

	_, err := bot.Send(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func ThreeData(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Data)
	msg.Text = "You press button 3"

	_, err := bot.Send(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func CreateCallbackRouter() *Router {
	router := NewRouter()

	router.register("1", OneData)
	router.register("2", TwoData)
	router.register("3", ThreeData)

	return router
}
