package main

import (
	// "fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("5702967487:AAFeA8A0gNKuYcZF2jlmmsSigd4i6tvt_To")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message

			command := strings.Split(update.Message.Text, " ")

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, command[0])
			// msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

//
//
