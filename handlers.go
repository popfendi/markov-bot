package main

import (
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/popfendi/markov-bot/markov"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var groupId int64 = parseEnvToI64(os.Getenv("GROUP_ID"))
var admin1Id int64 = parseEnvToI64(os.Getenv("ADMIN1_ID"))
var admin2Id int64 = parseEnvToI64(os.Getenv("ADMIN2_ID"))
var myId int64 = parseEnvToI64(os.Getenv("MY_ID"))
var minInterval = 300 // min interval between auto msgs
var maxInterval = 360 // max interval between auto msgs
var speakFreely = true

func updateHandler(update tgbotapi.Update) {
	isCommand := update.Message != nil && update.Message.IsCommand()
	isNonCommand := update.Message != nil && !update.Message.IsCommand()

	if isCommand {
		commandHandler(update)
	} else if isNonCommand {
		nonCommandHandler(update)
	}
}

func nonCommandHandler(update tgbotapi.Update) {
	isReplyToBot := update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == myId

	if isReplyToBot {
		replyHandler(update)
	} else {
		trainHandler(update)
	}
}

func commandHandler(update tgbotapi.Update) {

	switch update.Message.Command() {
	case "start":

	case "speak":
		speakHandler(update)
	case "speakinterval":
		setSpeakIntervalHandler(update)
	case "speakfreely":
		toggleSpeakFreelyHandler(update)
	}
}

func setSpeakIntervalHandler(update tgbotapi.Update) {
	if update.Message.From.ID == admin1Id || update.Message.From.ID == admin2Id {
		var re = regexp.MustCompile(`(?m)/speakinterval\s*(\d*)\s*(\d*)`)
		matches := re.FindAllStringSubmatch(update.Message.Text, -1)

		min, minErr := strconv.Atoi(matches[0][1])
		max, maxErr := strconv.Atoi(matches[0][2])

		if minErr == nil && maxErr == nil {
			maxInterval = max
			minInterval = min
		}
	}
}

func toggleSpeakFreelyHandler(update tgbotapi.Update) {
	if update.Message.From.ID == admin1Id || update.Message.From.ID == admin2Id {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if speakFreely {
			msg.Text = "I'll be quiet now."
			speakFreely = !speakFreely
		} else {
			msg.Text = "I'll say what's on my mind."
			speakFreely = !speakFreely
		}

		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := Bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}

func speakHandler(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyToMessageID = update.Message.MessageID
	if update.Message.Chat.ID == groupId {
		msg.Text = markov.Generate()
	} else {
		msg.Text = "Come and join the official Group!"
	}

	if _, err := Bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func trainHandler(update tgbotapi.Update) {
	if update.Message.Chat.ID == groupId && !update.Message.From.IsBot {
		log.Println("Training in: " + update.Message.Text)
		markov.Train(update.Message.Text)

		if mentionedMe(update.Message.Text) {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = markov.Generate()
			msg.ReplyToMessageID = update.Message.MessageID

			if _, err := Bot.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}
}

func mentionedMe(message string) bool {
	m := strings.ToLower(message)

	return strings.Contains(m, "@"+os.Getenv("BOT_USERNAME"))
}

func replyHandler(update tgbotapi.Update) {
	if update.Message.Chat.ID == groupId && !update.Message.From.IsBot {
		log.Println("Training in: " + update.Message.Text)
		markov.Train(update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.Text = markov.Generate()
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := Bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}

func speakOnInterval() {
	if speakFreely {
		for {
			msg := tgbotapi.NewMessage(groupId, "")
			msg.Text = markov.Generate()

			if _, err := Bot.Send(msg); err != nil {
				log.Println(err)
			}

			num := rand.Intn(maxInterval-minInterval) + minInterval
			time.Sleep(time.Duration(num) * time.Second)
		}
	}
}

func parseEnvToI64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println(err)
	}

	return i
}
