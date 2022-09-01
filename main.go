package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type binanceResp struct {
	Price float64 `json:"price,string"`
	Code  int64   `json:"code"`
}

type wallet map[string]float64

var dataBase = map[int64]wallet{}

func main() {
	bot, err := tgbotapi.NewBotAPI("5702967487:AAFeA8A0gNKuYcZF2jlmmsSigd4i6tvt_To")
	if err != nil {
		log.Panic(err)
	}
	log.Printf("%s ready to work", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		commandFromCustomer := strings.Split(update.Message.Text, " ")
		log.Println(update.Message.Text)

		switch commandFromCustomer[0] {
		case "Add":
			currencyToUpperCase := strings.ToUpper(commandFromCustomer[1])
			if len(commandFromCustomer) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не правильный формат команды(к примеру должно быть: Add btc 0.15)"))
				continue
			}
			summ, err := strconv.ParseFloat(commandFromCustomer[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Невозможно сконвертировать сумму"))
				continue
			}
			if _, ok := dataBase[update.Message.Chat.ID]; !ok {
				dataBase[update.Message.Chat.ID] = wallet{}
			}
			data, err := getPrice(currencyToUpperCase)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Указан неизвестный формат либо аббревиатура валюты"))
				continue
			}
			if data != 0 {
				dataBase[update.Message.Chat.ID][currencyToUpperCase] += summ
			}
			balanceText := fmt.Sprintf("Баланс: %s %.4f", currencyToUpperCase, dataBase[update.Message.Chat.ID][currencyToUpperCase])
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))
		case "Sub":
			currencyToUpperCase := strings.ToUpper(commandFromCustomer[1])
			if len(commandFromCustomer) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не правильный формат команды(к примеру должно быть: Add btc 0.15)"))
				continue
			}
			summ, err := strconv.ParseFloat(commandFromCustomer[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Невозможно сконвертировать сумму"))
				continue
			}
			if _, ok := dataBase[update.Message.Chat.ID]; !ok {
				dataBase[update.Message.Chat.ID] = wallet{}
			}
			if dataBase[update.Message.Chat.ID][currencyToUpperCase]-summ < 0 {
				msgZeroBalance := fmt.Sprintf("Неправильная операция по валюте %s в вашем кошельке всего %f, a вы хотите отнять %f", currencyToUpperCase, dataBase[update.Message.Chat.ID][currencyToUpperCase], summ)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msgZeroBalance))
				continue
			}
			if dataBase[update.Message.Chat.ID][currencyToUpperCase]-summ >= 0 {
				dataBase[update.Message.Chat.ID][currencyToUpperCase] -= summ
				balanceText := fmt.Sprintf("Баланс: %s %.4f", currencyToUpperCase, dataBase[update.Message.Chat.ID][currencyToUpperCase])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))
			}
		case "Show":
			msg := "Баланс: \n"
			var usdAmount float64
			for k, v := range dataBase[update.Message.Chat.ID] {
				coinPrice, err := getPrice(k)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				}
				usdAmount += v * coinPrice
				msg += fmt.Sprintf(" %s: %f [in $: %.2f]\n", k, v, v*coinPrice)
			}
			msg += fmt.Sprintf("Всего USD в кошельке: %.2f\n", usdAmount)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		case "Del":
			currencyToUpperCase := strings.ToUpper(commandFromCustomer[1])
			delete(dataBase[update.Message.Chat.ID], currencyToUpperCase)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Валюта удалена"))
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
		case "Отменить":
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
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Я пока еще не знаю такой команды"))
		}
	}
	if err != nil {
		fmt.Println(err)
	}
}

func getPrice(coin string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", coin))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var jsonResp binanceResp
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return 0, err
	}

	if jsonResp.Code != 0 {
		err = errors.New("Не корректная валюта")
		return 0, err
	}

	price := jsonResp.Price

	return price, nil
}
