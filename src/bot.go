package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var AppBot *tgbotapi.BotAPI

func bot_start() {
	AppBot, err := tgbotapi.NewBotAPI(Config.TgApi)
	command_router := CreateCommandRouter()
	callback_router := CreateCallbackRouter()

	if err != nil {
		log.Panicln(err)
	}

	AppBot.Debug = true

	botinfo, err := AppBot.GetMe()

	if err != nil {
		log.Panicln(err)
	}

	log.Printf("Auth %v token: %v", botinfo.UserName, Config.TgApi)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := AppBot.GetUpdatesChan(updateConfig)

	for update := range updates {
		//	log.Println(update)

		if update.Message != nil {
			if update.Message.Text[0] == '/' {
				command_router.run(update.Message.Command(), update, AppBot)
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			if _, err := AppBot.Request(callback); err != nil {
				log.Panicln(err)
			}

			callback_router.run(update.CallbackQuery.Data, update, AppBot)
		}

	}
}

func SomeText(update tgbotapi.Update) {

}
