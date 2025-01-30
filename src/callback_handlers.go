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
	msg.Text = "Enter name"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Cancel ❌", fmt.Sprintf("cancel_%v", user.CreateHabit)),
		),
	)

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

	msg_Text := fmt.Sprintf("Your habits %v", update.CallbackQuery.From.UserName)
	var msg_ReplyMarkup tgbotapi.InlineKeyboardMarkup

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

		return_button := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Return ↩️", "myhabitsreturn"),
		)

		habit_buttons = append(habit_buttons, return_button)

		msg_ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(habit_buttons...)
	} else {
		msg_Text = "You dont have habits("
		return_button := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Return ↩️", "myhabitsreturn"),
		)

		habit_buttons = append(habit_buttons, return_button)
		msg_ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(habit_buttons...)
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID,
		update.CallbackQuery.Message.MessageID,
		msg_Text,
		msg_ReplyMarkup,
	)

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
	msg_text := fmt.Sprintf("Name: %v\nTime: %v", habit.Name, habit.RemTime)
	msg_ReplyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Check %v", is_check), fmt.Sprintf("complete_%v", habit.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Delete ❌", fmt.Sprintf("delete_%v", habit.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Return ↩️", "habitreturn"),
		),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID,
		update.CallbackQuery.Message.MessageID,
		msg_text,
		msg_ReplyMarkup,
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
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Return ↩️", "habitreturn"),
			),
		)

		msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, markup)

		bot.Send(msg)
	}
}

func CancelCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	args := strings.Split(data, "_")

	habitid, _ := strconv.ParseInt(args[1], 10, 64)
	habit := Habit{}
	habit.Load(habitid, db)
	habit.Delete(db)

	msg := tgbotapi.NewDeleteMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID)

	_, err := bot.Request(msg)

	if err != nil {
		log.Panicln(err)
	}
}

func HabitDeleteCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	args := strings.Split(data, "_")

	habitid, _ := strconv.ParseInt(args[1], 10, 64)

	msg_Text := "Confirm delete"
	msg_ReplyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Confirm ✅", fmt.Sprintf("confirmdelete_%v", habitid)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Cancel ❌", fmt.Sprintf("canceldelete_%v", habitid)),
		),
	)
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.From.ID,
		update.CallbackQuery.Message.MessageID,
		msg_Text,
		msg_ReplyMarkup,
	)

	bot.Send(msg)
}

func ConfirmDeleteCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	args := strings.Split(data, "_")

	habitid, _ := strconv.ParseInt(args[1], 10, 64)
	habit := Habit{}
	habit.Load(habitid, db)
	habit.Delete(db)

	HabitReturnCallbackHandler(data, update, bot, db)
}

func CancelDeleteCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	HabitMenuCallbackHandler(data, update, bot, db)
}

func StatisticCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {

}

func MyHabitsReturnCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	StartCommand("return", update, bot, db)
}

func HabitReturnCallbackHandler(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	MyHabitsCallbackHandler(data, update, bot, db)
}

func CreateCallbackRouter() *Router {
	router := NewRouter()

	router.register("create", CreateCallbackHandler)
	router.register("myhabits", MyHabitsCallbackHandler)
	router.register("statistic", StatisticCallbackHandler)
	router.register("habit", HabitMenuCallbackHandler)
	router.register("complete", HabitCheckCallbackHandler)
	router.register("delete", HabitDeleteCallbackHandler)
	router.register("habitreturn", HabitReturnCallbackHandler)
	router.register("myhabitsreturn", MyHabitsReturnCallbackHandler)
	router.register("cancel", CancelCallbackHandler)
	router.register("confirmdelete", ConfirmDeleteCallbackHandler)
	router.register("canceldelete", CancelDeleteCallbackHandler)

	return router
}
