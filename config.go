package main

import (
	"log"
	"os"
)

// ConfigData contains config data. Api keys, paths
type ConfigData struct {
	botAPIKey     string
	weatherAPIKey string
	dbPath        string
}

func loadConfigFromEnv() ConfigData {

	botAPIKey, isSet := os.LookupEnv("TELEGRAM_BOT_API_KEY")
	if !isSet {
		log.Panic("No bot api key found. Please set TELEGRAM_BOT_API_KEY env")
	}

	weatherAPIKey, isSet := os.LookupEnv("WEATHER_BOT_API")
	if !isSet {
		log.Panic("No bot api key found. Please set WEATHER_BOT_API env")
	}
	dbPath, isSet := os.LookupEnv("TELEGRAM_BOT_DB_PATH")
	if !isSet {
		log.Panic("No bot api key found. Please set TELEGRAM_BOT_DB_PATH env")
	}

	log.Println(botAPIKey + " " + weatherAPIKey + " " + dbPath)

	return ConfigData{botAPIKey, weatherAPIKey, dbPath}
}
