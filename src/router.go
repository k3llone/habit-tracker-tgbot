package main

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handle func(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB)
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

func (r *Router) run(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) bool {
	handler, exists := r.Map[data]

	if exists {
		handler(update, bot, db)
	}

	return exists
}

type TextRouter struct {
	Map map[string]TextHandle
}

func NewTextRouter() *TextRouter {
	return &TextRouter{make(map[string]TextHandle)}
}

func (r *TextRouter) register(state string, handle TextHandle) {
	r.Map[state] = handle
}

func (r *TextRouter) run(text string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) bool {
	user := User{}
	err := user.Load(update.Message.From.ID, db)
	log.Println(err)

	state := user.State[len(user.State)-1]

	handler, exists := r.Map[state]

	log.Println(user)

	log.Printf("Text handler %v %v exists: %v", state, text, exists)

	if exists {
		handler(text, update, bot, db)
	}

	return exists
}
