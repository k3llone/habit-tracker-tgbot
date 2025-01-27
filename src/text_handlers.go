package main

import (
	"database/sql"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartTextHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	msg := tgbotapi.NewMessage(update.Message.From.ID, "")
	msg.Text = "Menu state"

	_, err := bot.Send(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func CreateNameTextHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	user := User{}
	habit := Habit{}
	user.Load(update.Message.From.ID, db)

	habit.Name = data
	habit.UserId = user.Id

	habit.Insert(db)

	user.CreateHabit = habit.Id

	user.State = "create_time"
	user.Update(db)

	msg := tgbotapi.NewMessage(update.Message.From.ID, "Enter time (hour:minute)")

	_, err := bot.Send(msg)

	if err != nil {
		log.Println(err)
	}
}

func CreateTimeTextHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	user := User{}
	habit := Habit{}
	user.Load(update.Message.From.ID, db)
	habit.Load(user.CreateHabit, db)

	user.State = "menu"
	user.Update(db)

	habit.RemTime = data

	habit.Update(db)

	msg := tgbotapi.NewMessage(update.Message.From.ID, "")
	msg.Text = fmt.Sprintf("New habit created:\nName: %v\nTime: %v", habit.Name, habit.RemTime)

	_, err := bot.Send(msg)

	if err != nil {
		log.Println(err)
	}
}

func CreateTextRouter() *Router {
	router := NewRouter()

	router.register("menu", StartTextHandler)
	router.register("create_name", CreateNameTextHandler)
	router.register("create_time", CreateTimeTextHandler)

	return router
}
