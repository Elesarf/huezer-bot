package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/proxy"
)

// SystemConfig contains system config
var SystemConfig ConfigData

// ProxyAwareHTTPClient for using proxy
func ProxyAwareHTTPClient() *http.Client {
	var dialer proxy.Dialer
	dialer = proxy.Direct
	// read env and, if set proxy, apply
	proxyServer, isSet := os.LookupEnv("HTTP_PROXY")
	if isSet {
		proxyURL, err := url.Parse(proxyServer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid proxy url %q\n", proxyURL)
		}
		dialer, err = proxy.FromURL(proxyURL, proxy.Direct)
		_check(err)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	return httpClient
}

func _check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	log.Printf("Start")

	SystemConfig = loadConfigFromEnv()

	client := ProxyAwareHTTPClient()
	// create bot using token, client
	bot, err := tgbotapi.NewBotAPIWithClient(SystemConfig.botAPIKey, client)
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
