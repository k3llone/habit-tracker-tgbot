package main

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func bot_start(db *sql.DB) {
	AppBot, err := tgbotapi.NewBotAPI(Config.TgApi)

	command_router := CreateCommandRouter()
	callback_router := CreateCallbackRouter()
	text_router := CreateTextRouter()

	if err != nil {
		log.Panicln(err)
	}

	AppBot.Debug = false

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
				command_router.RunCommand(update.Message.Command(), update, AppBot, db)
			} else {
				text_router.RunText(update.Message.Text, update, AppBot, db)
			}

		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

			if _, err := AppBot.Request(callback); err != nil {
				log.Panicln(err)
			}

			callback_router.RunCallback(update.CallbackQuery.Data, update, AppBot, db)
		}

	}
}

func SomeText(update tgbotapi.Update) {

}
