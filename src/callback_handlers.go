package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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

func MyHabitsCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	user := User{}
	user.Load(update.CallbackQuery.From.ID, db)
	user.State = "myhabits"
	user.Update(db)

	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "")
	msg.Text = fmt.Sprintf("Your habits %v", update.CallbackQuery.From.UserName)

	habit_buttons := [][]tgbotapi.InlineKeyboardButton{}
	habits := make([]Habit, 0)

	res, err := db.Query("SELECT * FROM Habits WHERE user=$1", user.Id)

	if err != nil {
		log.Panicln(err)
	}

	for res.Next() {
		h := Habit{}

		res.Scan(&h.Id, &h.UserId, &h.Name, &h.RemTime)

		habits = append(habits, h)
	}

	if len(habits) != 0 {
		for _, v := range habits {

			is_check := ""

			if CheckHabitToday(v.UserId, v.Id, db) {
				is_check = "✅"
			} else {
				is_check = "❌"
			}

			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v %v %v", v.Name, v.RemTime, is_check), fmt.Sprintf("habit_%v", v.Id)),
			)

			habit_buttons = append(habit_buttons, row)
		}

		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(habit_buttons...)
	} else {
		msg.Text = "You dont have habits("
	}

	_, err = bot.Send(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func HabitMenuCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	args := strings.Split(data, "_")

	user := User{}
	user.Load(update.CallbackQuery.From.ID, db)

	habitid, _ := strconv.ParseInt(args[1], 10, 64)

	habit := Habit{}
	habit.Load(habitid, db)

	is_check := ""

	if CheckHabitToday(habit.UserId, habitid, db) {
		is_check = "✅"
	} else {
		is_check = "❌"
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "")
	msg.Text = fmt.Sprintf("Name: %v\nTime: %v", habit.Name, habit.RemTime)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Check %v", is_check), fmt.Sprintf("complete_%v", habit.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Delete ❌", fmt.Sprintf("delete_%v", habit.Id)),
		),
	)

	if habit.UserId == user.Id {
		_, err := bot.Send(msg)

		if err != nil {
			log.Panicln(err)
		}
	}
}

func HabitCheckCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	args := strings.Split(data, "_")

	habitid, _ := strconv.ParseInt(args[1], 10, 64)
	habit := Habit{}
	habit.Load(habitid, db)

	if !CheckHabitToday(habit.UserId, habit.Id, db) {
		habit_check := HabitComplete{}
		habit_check.HabitId = habitid

		today := time.Now()
		check_time := fmt.Sprintf("%v.%v.%v", today.Day(), today.Month(), today.Year())

		habit_check.Date = check_time
		err := habit_check.Insert(db)

		if err != nil {
			log.Panicln(err)
		}

		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Check ✅", fmt.Sprintf("complete_%v", habit.Id)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Delete ❌", fmt.Sprintf("delete_%v", habit.Id)),
			),
		)

		msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, markup)

		bot.Send(msg)
	}
}

func HabitDeleteCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	args := strings.Split(data, "_")

	habitid, _ := strconv.ParseInt(args[1], 10, 64)
	habit := Habit{}
	habit.Load(habitid, db)
	habit.Delete(db)
}

func StatisticCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {

}

func CreateCallbackRouter() *Router {
	router := NewRouter()

	router.register("create", CreateCallbackHandler)
	router.register("myhabits", MyHabitsCallbackHandler)
	router.register("statistic", StatisticCallbackHandler)
	router.register("habit", HabitMenuCallbackHandler)
	router.register("complete", HabitCheckCallbackHandler)
	router.register("delete", HabitDeleteCallbackHandler)

	return router
}
