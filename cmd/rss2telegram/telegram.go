package main

import (
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	tgBot *tgbotapi.BotAPI
	chat  int64
)

func initTelegram(token string, chatID int64) {
	chat = chatID

	var err error

	tgBot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize telegram bot: %w", err))
	}
}

func serveTelegram(in <-chan string) {
	for msg := range in {
		_, err := tgBot.Send(tgbotapi.NewMessage(chat, msg))
		if err != nil {
			log.Println(fmt.Errorf("failed to send telegram link %s on %d: %w", msg, chat, err))
		}
	}
}
