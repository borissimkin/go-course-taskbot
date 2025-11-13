package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramDelivery struct {
	bot *tgbotapi.BotAPI
}

func (d *TelegramDelivery) Send(chatID ID, text string) {
	msg := tgbotapi.NewMessage(int64(chatID), text)

	_, err := d.bot.Send(msg)

	if err != nil {
		log.Printf("error send to chatID=%d text=%s err=%s", chatID, text, err)
	}
}

func (d *TelegramDelivery) Receive(body io.Reader) (RequestMessage, error) {
	var message tgbotapi.Update

	err := json.NewDecoder(body).Decode(&message)
	if err != nil {
		return RequestMessage{}, fmt.Errorf("json decode error in TelegramDelivery.Receive: %w", err)
	}

	return RequestMessage{
		User: User{
			ID:       ID(message.Message.From.ID),
			UserName: message.Message.From.UserName,
		},
		ChatID: ID(message.Message.Chat.ID),
		Text:   message.Message.Text,
	}, nil
}
