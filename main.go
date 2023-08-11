package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"telegram-weather-bot/handlers"
	"telegram-weather-bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	var cfg models.Config

	// Чтение конфигурации из файла config.yaml
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Fatalf("cannot load configuration: %v", err)
	}

	// Инициализация Telegram бота
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatalf("Failed to create new Telegram bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Устанавливаем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to get updates channel: %v", err)
	}

	// Инициализация канала для graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Основной цикл для обработки сообщений
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			// Обработка входящих сообщений и команд
			handlers.HandleMessage(bot, &cfg, update.Message)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	<-shutdown

	// Закрытие канала обновлений и завершение работы
	bot.StopReceivingUpdates()
	log.Println("Bot stopped gracefully")
}
