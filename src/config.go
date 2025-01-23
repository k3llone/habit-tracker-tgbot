package main

import (
	"encoding/json"
	"log"
	"os"
)

type AppCofnig struct {
	TgApi string `json:"tgapi"`
}

var Config AppCofnig

func load_config() {
	content, err := os.ReadFile("config.json")

	if err != nil {
		log.Panicln(err)
	}

	err = json.Unmarshal(content, &Config)

	if err != nil {
		log.Panicln(err)
	}
}
