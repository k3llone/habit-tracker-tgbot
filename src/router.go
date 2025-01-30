package main

import (
	"database/sql"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handle func(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB)
type TextHandle func(text string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB)

type Router struct {
	Map map[string]Handle
}

func NewRouter() *Router {
	return &Router{make(map[string]Handle)}
}

func (r *Router) register(data string, handle Handle) {
	r.Map[data] = handle
}

func (r *Router) RunCommand(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) bool {
	handler, exists := r.Map[data]

	if exists {
		handler(data, update, bot, db)
	}

	return exists
}

func (r *Router) RunCallback(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) bool {
	args := strings.Split(data, "_")
	handler, exists := r.Map[args[0]]

	if exists {
		handler(data, update, bot, db)
	}

	return exists
}

func (r *Router) RunText(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) bool {
	user := User{}
	user.Load(update.Message.From.ID, db)
	//	log.Println(err)

	handler, exists := r.Map[user.State]

	//	log.Println(user)

	//	log.Printf("Text handler %v %v exists: %v", user.State, data, exists)

	if exists {
		handler(data, update, bot, db)
	}

	return exists
}
