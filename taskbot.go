package main

import (
	"context"
	"log"
	http "net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	WebhookURL = "http://127.0.0.1:8081"
	BotToken   = "_golangcourse_test"
)

func startTaskBot(ctx context.Context, httpListenAddr string) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}

	tgDelivery := &TelegramDelivery{
		bot: bot,
	}

	repo := NewRepo()

	taskTracker := &TaskTracker{
		delivery: tgDelivery,
		repo:     repo,
	}

	router := &CmdRouter{
		service: taskTracker,
	}

	srv := Server{
		receiver: tgDelivery,
		router: router,
	}

	bot.Debug = true

	http.ListenAndServe(httpListenAddr, srv)

	// сюда писать код
	/*
		в этом месте вы стартуете бота,
		стартуете хттп сервер который будет обслуживать этого бота
		инициализируете ваше приложение
		и потом будете обрабатывать входящие сообщения
	*/
	return nil
}

func main() {
	err := startTaskBot(context.Background(), ":8081")
	if err != nil {
		log.Fatalln(err)
	}
}

// это заглушка чтобы импорт сохранился
func __dummy() {
	tgbotapi.APIEndpoint = "_dummy"
}
