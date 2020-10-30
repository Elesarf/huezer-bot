package main

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

// SystemConfig contains system config
var SystemConfig ConfigData

func _check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	log.Printf("Start")

	SystemConfig = loadConfigFromEnv()

	// create bot using token, client
	bot, err := tgbotapi.NewBotAPI(SystemConfig.botAPIKey)
	_check(err)

	// debug mode on
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// set update interval
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1000

	updates, err := bot.GetUpdatesChan(u)
	_check(err)

	var message InternalMessage
	// get new updates
	for update := range updates {

		if update.CallbackQuery != nil {

			var config tgbotapi.CallbackConfig
			config.URL = "https://huezer.xyz"
			config.CallbackQueryID = update.CallbackQuery.ID

			log.Println("--------------------")
			log.Println(config)
			log.Println("--------------------")
			bot.AnswerCallbackQuery(config)
		}

		// if message from channel
		if update.ChannelPost != nil {
			message = InternalMessage{update.ChannelPost.Chat.ID,
				0,
				"",
				strings.ToLower(update.ChannelPost.Text)}
		}

		// if message from user
		if update.Message != nil {
			message = InternalMessage{update.Message.Chat.ID,
				update.Message.MessageID,
				update.Message.From.UserName,
				strings.ToLower(update.Message.Text)}
		}
		processMessage(&message, bot)
	}
}

func processMessage(message *InternalMessage, bot *tgbotapi.BotAPI) {
	messageParse(message, bot)
	if message.messageText != "" {
		sendAnswer(message, bot)
	}
}

func sendAnswer(message *InternalMessage, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.chatID, message.messageText)
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}
