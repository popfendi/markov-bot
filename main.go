package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/popfendi/markov-bot/markov"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func main() {
	markov.Init()

	b, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	Bot = b
	Bot.Debug = true

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	go speakOnInterval()

	for update := range updates {
		log.Printf("New Update [%s]", fmt.Sprint(update.UpdateID))

		jsonObj, _ := json.Marshal(update)
		fmt.Println(string(jsonObj))

		updateHandler(update)

	}

}
