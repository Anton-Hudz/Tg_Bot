package main

import (
	// "errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// var (
// 	problemWithParseErr = errors.New("неверная команда")
// )

type wallet map[string]float64

var dataBase = map[int64]wallet{}

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
			// if update.Message == nil {
			// 	continue
			// }
			command := strings.Split(update.Message.Text, " ")

			switch command[0] {
			case "Add":
				if len(command) != 3 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не правильный формат команды(к примеру должно быть: Add btc 0.15) "))
				}
				amount, err := strconv.ParseFloat(command[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				}
				if _, ok := dataBase[update.Message.Chat.ID]; !ok {
					dataBase[update.Message.Chat.ID] = wallet{}
				}
				dataBase[update.Message.Chat.ID][command[1]] += amount

				balanceText := fmt.Sprintf("%f", dataBase[update.Message.Chat.ID][command[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))
			case "Show":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ну ты и коза"))

			case "Sub":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Сама "))
			case "Привет":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Привет!\n Я тестовый бот, умею показывать интересные новости, отправлять сообщения твоим контактам и делится с ними интересными новостями \nнапиши мне 'Покажи', чтобы увидеть новую интересную новость..."))
			case "Покажи":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "https://ru.wikihow.com/признаться-другу,-что-вы-гей-и-что-вы-любите-его"))
				time.Sleep((time.Second * 2))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Чтобы отправить эту новость 10 твоим самым часто используемым контактам введи 'Отправить', чтобы отменить отправку введи 'Отмена', в случае выхода из бота новость будет отправлена автоматически"))
			case "Отправить":
				time.Sleep((time.Second * 1))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Отправляю последнюю новость 10 твоим самым часто используемым контактам"))
				time.Sleep((time.Second * 1))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "3..."))
				time.Sleep((time.Second * 1))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "2..."))
				time.Sleep((time.Second * 1))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "1..."))
				time.Sleep((time.Second * 1))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Сообщение отправлено"))
			case "Отмена":
				time.Sleep((time.Second * 1))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Отправляю последнюю новость 10 твоим самым часто используемым контактам"))
				time.Sleep((time.Second * 1))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "3..."))
				time.Sleep((time.Second * 1))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "2..."))
				time.Sleep((time.Second * 1))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "1..."))
				time.Sleep((time.Second * 1))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Сообщение отправлено"))

			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "я пока еще не знаю такой команды"))
			}

			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, command[0])
			// // msg.ReplyToMessageID = update.Message.MessageID

			// bot.Send(msg)

		}
	}
	if err != nil {
		fmt.Println(err)
	}
}

// https://pikabu.ru/story/o_tom_kak_ya_osoznal_sebya_ili_na_grani_propasti_6354211
