package main

import (
	"log"
	"tg-quiz/internal/bot"
	"tg-quiz/internal/game"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	adminID = 1234567890 // Замените на ваш Telegram ID
)

func main() {
	// Инициализация бота
	botAPI, err := tgbotapi.NewBotAPI("7989564521:AAGsgERw1ETO_NEm0YlTOk3NO31ESYyfUB0")
	if err != nil {
		log.Panic(err)
	}

	botAPI.Debug = true
	log.Printf("Бот %s успешно запущен", botAPI.Self.UserName)

	// Создаем игру и обработчик
	game := game.NewGame(adminID)
	handler := bot.NewHandler(botAPI, game)

	// Настройка получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botAPI.GetUpdatesChan(u)

	// Обработка сообщений
	for update := range updates {
		if update.Message != nil {
			handler.HandleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			handler.HandleCallback(update.CallbackQuery)
		}
	}
}
