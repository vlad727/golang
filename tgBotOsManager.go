package main

import (
	"log"
	"os"
	"runtime"
	

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// function get hostname
func getHostname() string {
	hostname, err := os.Hostname()

	if err != nil {
		log.Printf("Can not get hostname name")
	}
	return hostname
}

// function get type OS
func typeOS() string {
	os := runtime.GOOS
	return os

}
func main() {
	tokenString := os.Getenv("TELEGRAM_APITOKEN")
	bot, err := tgbotapi.NewBotAPI(tokenString)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			switch update.Message.Text {
			case "/help":
				res := "You can get info about server:\n/hostname\n/typeos\n"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			case "/hostname":
				res := getHostname()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			case "/typeos":
				res := typeOS()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)

			}

		}
	}
}
