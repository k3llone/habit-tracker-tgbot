package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handle func(update tgbotapi.Update, bot *tgbotapi.BotAPI)

type Router struct {
	Map map[string]Handle
}

func NewRouter() *Router {
	return &Router{make(map[string]Handle)}
}

func (r *Router) register(data string, handle Handle) {
	r.Map[data] = handle
}

func (r *Router) run(data string, update tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	handler, exists := r.Map[data]

	if exists {
		handler(update, bot)
	}

	return exists
}
