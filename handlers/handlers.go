package handlers

import (
	"net/http"
	"telegram-weather-bot/models"
	"telegram-weather-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleMessage обрабатывает входящие сообщения и команды от пользователей.
func HandleMessage(bot *tgbotapi.BotAPI, cfg *models.Config, message *tgbotapi.Message) {
	if message.IsCommand() {
		handleCommand(bot, cfg, message)
	} else {
		handleLocationRequest(bot, cfg, message)
	}
}

// handleCommand обрабатывает команды от пользователей.
func handleCommand(bot *tgbotapi.BotAPI, cfg *models.Config, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(message.Chat.ID, "Привет! Отправь мне название города или поделись своей геолокацией, чтобы узнать погоду.")
		msg.ReplyMarkup = getKeyboard()
		bot.Send(msg)
	case "hide":
		hideKeyboard(bot, message.Chat.ID)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Извините, я не понимаю эту команду.")
		bot.Send(msg)
	}
}

// handleLocationRequest обрабатывает запросы геолокации от пользователей.
func handleLocationRequest(bot *tgbotapi.BotAPI, cfg *models.Config, message *tgbotapi.Message) {
	if message.Location != nil {
		// Здесь мы будем обрабатывать геолокацию и делать запрос к API погоды
		weatherInfo, err := services.GetWeatherByLocation(message.Location.Latitude, message.Location.Longitude, cfg.WeatherAPIKey)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении погоды. Попробуйте позже.")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, weatherInfo)
		bot.Send(msg)
	} else {
		weatherInfo, err := services.GetWeatherByCityName(http.DefaultClient, message.Text, cfg.WeatherAPIKey)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении погоды для указанного города. Попробуйте позже.")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, weatherInfo)
		bot.Send(msg)
	}
}

// клавиатура с кнопками "Отправить геопозицию" и "Помощь"
func getKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation("Отправить геопозицию"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Помощь"),
		),
	)
	return keyboard
}

// hideKeyboard скрывает ReplyKeyboard для указанного chatID.
func hideKeyboard(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Клавиатура скрыта.")
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}
	bot.Send(msg)
}
